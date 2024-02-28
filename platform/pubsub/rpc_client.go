package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResponseMessage struct {
	Result    []byte
	MessageID string
}

type PubsubClient struct {
	projectID            string
	topicName            string
	subscriptionName     string
	callbacks            map[string]func([]byte) ([]byte, error)
	responseQueue        chan *ResponseMessage
	internalLock         sync.Mutex
	callbackExecComplete chan struct{}
	stopEvent            chan struct{}
	message_counter      int
	client               *pubsub.Client
	topic                *pubsub.Topic
	sub                  *pubsub.Subscription
	sendMessageIDs       map[string]interface{}
}

func NewPubSubClient(projectID, topicName, subName string) (*PubsubClient, error) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, os.Getenv("GCS_PROJECT_ID"))
	if err != nil {
		return nil, fmt.Errorf("failed to create Pub/Sub client : %v", err)
	}

	topic := client.Topic(topicName)
	subscription := client.Subscription(subName)

	exists, err := subscription.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		_, err = client.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 20 * time.Second,
		})
		if err != nil {
			return nil, err
		}
	}

	if err != nil && status.Code(err) != codes.AlreadyExists {
		return nil, fmt.Errorf("failed to create Pub/Sub subscription: %v", err)
	}

	return &PubsubClient{
		projectID:            projectID,
		topicName:            topicName,
		subscriptionName:     subName,
		client:               client,
		topic:                topic,
		sub:                  subscription,
		callbacks:            make(map[string]func([]byte) ([]byte, error)),
		responseQueue:        make(chan *ResponseMessage),
		internalLock:         sync.Mutex{},
		callbackExecComplete: make(chan struct{}),
		stopEvent:            make(chan struct{}),
		message_counter:      0,
		sendMessageIDs:       make(map[string]interface{}),
	}, nil
}

func (c *PubsubClient) StartListening() {
	go func() {
		for {
			select {
			case <-c.stopEvent:
				return
			default:
				_, err := c.ListenForMessages()
				if err != nil {
					log.Printf("Error while listening for messages: %v\n", err)
				}
			}
		}
	}()
}

func (c *PubsubClient) StopListening() {
	close(c.stopEvent)
}

func (c *PubsubClient) ListenForMessages() (interface{}, error) {
	ctx := context.Background()

	var result interface{}

	err := c.sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		c.internalLock.Lock()
		defer c.internalLock.Unlock()

		if msg.Data == nil {
			log.Printf("Received an empty Pub/Sub message")
			msg.Ack()
			return
		}

		messageData := msg.Data
		fmt.Printf("Received Message : %s\n", messageData)
		var messageMap map[string]interface{}
		if err := json.Unmarshal(messageData, &messageMap); err != nil {
			log.Printf("Failed to decode JSON message: %v\n", err)
		} else {
			messageID := messageMap["message_id"].(string)
			if _, exists := c.sendMessageIDs[messageID]; exists {
				result = messageMap["result"]
				c.sendMessageIDs[messageID] = result
			}
		}

		msg.Ack()
	})

	return result, err
}

func (c *PubsubClient) PublishMessage(methodName string, args interface{}) (string, error) {
	c.internalLock.Lock()
	defer c.internalLock.Unlock()

	messageID := fmt.Sprintf("msg_%d", c.message_counter)
	messagePayload := map[string]interface{}{
		"method":     methodName,
		"args":       args,
		"kwargs":     "",
		"message_id": messageID,
	}

	messageData, err := json.Marshal(messagePayload)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	result := c.client.Topic(c.topicName).Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	_, err = result.Get(ctx)
	if err != nil {
		return "", err
	}

	c.message_counter++

	return messageID, nil
}

func (rpc *PubsubClient) Close() {
	rpc.client.Close()
}

func (c *PubsubClient) RemoteMethod(methodName string, callback func([]byte) ([]byte, error)) func() ([]byte, error) {
	c.callbacks[methodName] = callback

	return func() ([]byte, error) {
		c.internalLock.Lock()
		defer c.internalLock.Unlock()

		messageID, err := c.PublishMessage(methodName, nil)
		if err != nil {
			return nil, err
		}

		<-c.callbackExecComplete

		responseMessage := <-c.responseQueue
		if responseMessage != nil && responseMessage.MessageID == messageID {
			return responseMessage.Result, nil
		}

		log.Printf("Timeout waiting for response for message_id : %s\n", messageID)
		return nil, nil
	}
}

func (c *PubsubClient) WaitForResponse(messageID string, timeout time.Duration, deleteAfterUse bool) (interface{}, error) {
	c.internalLock.Lock()
	defer c.internalLock.Unlock()

	if result, exists := c.sendMessageIDs[messageID]; exists {
		if deleteAfterUse {
			delete(c.sendMessageIDs, messageID)
		}
		return result, nil
	}

	// if not found, setup a timer channel for the timeout
	timer := time.After(timeout)

	select {
	case <-timer:
		return nil, fmt.Errorf("timeout waiting for result for message_id: %s", messageID)
	case <-c.stopEvent:
		return nil, fmt.Errorf("RPC client stopped listening for messages")
	}
}
