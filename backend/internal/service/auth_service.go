package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidOTP       = errors.New("invalid or expired OTP")
	ErrUserNotVerified  = errors.New("user is not verified")
	ErrEmailAlreadyUsed = errors.New("email already in use")
)

const (
	otpLength      = 6
	otpExpiryMins  = 10
	jwtExpiryHours = 24
)

// GenerateOTP generates a secure random 6-digit OTP
func GenerateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// SendMockEmail simulates sending an email with the OTP
func SendMockEmail(email, otp string) {
	// In a real app, integrate with SendGrid, AWS SES, etc.
	log.Printf("==== MOCK EMAIL ====\nTo: %s\nSubject: Your Verification Code\nBody: Your OTP is: %s\n====================\n", email, otp)
}

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with a plain text password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token for a user
func GenerateJWT(userID string, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * jwtExpiryHours).Unix(),
	})

	return token.SignedString([]byte(secret))
}
