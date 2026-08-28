package models

import (
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	IsVerified   bool      `json:"is_verified"`
	CreatedAt    time.Time `json:"created_at"`
}

type OTP struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Image struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	OriginalFilename string    `json:"original_filename"`
	S3Key            string    `json:"s3_key"`
	Format           string    `json:"format"`
	Size             int64     `json:"size"`
	Width            int       `json:"width"`
	Height           int       `json:"height"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type Transformation struct {
	ID              string    `json:"id"`
	ImageID         string    `json:"image_id"`
	TransformParams string    `json:"transform_params"` // JSON string representation
	S3Key           string    `json:"s3_key"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}
