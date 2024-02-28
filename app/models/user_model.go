package models

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type UserUpdate struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AvatarURL string `json:"avatarUrl"`
}

type User struct {
	gorm.Model
	FirstName       string          `json:"firstName"`
	LastName        string          `json:"last_name"`
	Username        string          `json:"username"`
	AvatarURL       string          `json:"avatarUrl"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	Email           string          `json:"email" validate:"required,email,lte=255"`
	PasswordHash    string          `json:"password_hash,omitempty" validate:"required,lte=255"`
	UserStatus      int             `json:"user_status" validate:"required,len=1"`
	UserRole        string          `json:"user_role" validate:"required,lte=25"`
	UploadedImages  []UploadedImage `json:"uploaded_images" gorm:"type:json"`
	VisitedProducts []string        `json:"visited_products" gorm:"type:json"`
	VisitedPages    []string        `json:"visited_pages" gorm:"type:json"`
}

type UploadedImage struct {
	gorm.Model
	UserID      uint   `json:"user_id" gorm:"not null"`
	ContentType string `json:"content_type" gorm:"not null"`
	ContentURL  string `json:"content_url" gorm:"not null"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public" gorm:"default:false"`
	UploadDate  time.Time
}

type SearchHistory struct {
	gorm.Model
	Query      string    `json:"query"`
	SearchDate time.Time `json:"searchDate"`
	UserID     uint      `json:"-"`
}

func UploadAvatar(userID uint, avatarImage io.Reader) (string, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithAPIKey(os.Getenv("GCS_API_KEY")))
	if err != nil {
		return "", fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	// create a User ID basedd filename
	avatarFilename := fmt.Sprintf("%d_avatar.png", userID)
	bucket := client.Bucket(os.Getenv("AVATAR_BUCKET"))
	obj := bucket.Object(avatarFilename)

	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, avatarImage); err != nil {
		return "", fmt.Errorf("failed to write avatar image to GCS: %v", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer : %v", err)
	}

	// set public access to the avatar URL
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("failed to set GCS object ACL : %v", err)
	}

	avatarURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", os.Getenv("AVATAR_BUCKET"), avatarFilename)
	return avatarURL, nil
}

func UploadImageToCloudStorage(db *gorm.DB, userID uint, contentType, description string, isPublic bool, imageData []byte) (string, error) {
	ctx := context.Background()
	// Setup google cloud storage client
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create Google Cloud Storage Client : %v", err)
		return "", fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	// create a new uuid for image filename
	imageID := uuid.New()
	imageFilename := fmt.Sprintf("%s.jpg", imageID.String())

	// upload the image to GCS
	obj := client.Bucket(os.Getenv("GCS_IMAGE_BUCKET")).Object(imageFilename)
	wc := obj.NewWriter(ctx)
	defer wc.Close()

	wc.ContentType = contentType

	// write image data to the google cloud storage object
	if _, err := wc.Write(imageData); err != nil {
		log.Printf("Failed to write image data to Google Cloud Storage: %v", err)
		return "", fmt.Errorf("failed to write image data to GCS: %v", err)
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("error setting ACL: %v", err)
	}

	// construct the URL for the uploaded image
	imageURL := fmt.Sprintf("gs://%s/%s", os.Getenv("GCS_IMAGE_BUCKET"), imageFilename)

	uploadedImage := UploadedImage{
		UserID:      userID,
		ContentType: contentType,
		ContentURL:  imageURL,
		Description: description,
		IsPublic:    isPublic,
		UploadDate:  time.Now(),
	}

	result := db.Create(&uploadedImage)
	if result.Error != nil {
		log.Printf("Failed to save Upload Image record to the database: %v", result.Error)
		return "", fmt.Errorf("failed to save image record to database: %v", result.Error)
	}

	return imageURL, nil
}
