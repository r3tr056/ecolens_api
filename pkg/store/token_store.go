package store

import (
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type TokenStore struct {
	mu          sync.Mutex
	tokens      map[string]TokenInfo
	cron        *cron.Cron
	cronJobId   cron.EntryID
	cleanupSpec string
}

type TokenInfo struct {
	userId     string
	expiration time.Time
}

// NewTokenStore creates a new instace of TokenStore
func NewTokenStore(cleanUpSpec string) *TokenStore {
	ts := &TokenStore{
		tokens:      make(map[string]TokenInfo),
		cron:        cron.New(),
		cleanupSpec: cleanUpSpec,
	}

	ts.cronJobId, _ = ts.cron.AddFunc(ts.cleanupSpec, ts.CleanUpExpiredTokens)
	ts.cron.Start()

	return ts
}

func (ts *TokenStore) Add(token string, userID string, expiration time.Time) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.tokens[token] = TokenInfo{
		userId:     userID,
		expiration: expiration,
	}
}

func (ts *TokenStore) Get(token string) (userID string, expiration time.Time, found bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	info, ok := ts.tokens[userID]
	if !ok {
		return "", time.Time{}, false
	}

	return info.userId, info.expiration, true
}

func (ts *TokenStore) CleanUpExpiredTokens() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	currentTime := time.Now()
	for token, info := range ts.tokens {
		if info.expiration.Before(currentTime) {
			delete(ts.tokens, token)
		}
	}
}

func (ts *TokenStore) StopCleanupJob() {
	ts.cron.Remove(ts.cronJobId)
	ts.cron.Stop()
}
