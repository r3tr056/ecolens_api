package models

import "time"

type SignUp struct {
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
	UserRole string `json:"user_role" validate:"required,lte=255"`
}

type SignIn struct {
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
}

type UserMeta struct {
	UserID uint `json:"userId"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email,lte=255"`
}

type ResetTokenInfo struct {
	UserID         uint
	ExpirationTime time.Time
}
