package db

import (
	"context"
	"database/sql"
	"errors"

	"image_processing_backend/internal/models"
)

// CreateUser inserts a new user and returns the inserted ID
func (s *Store) CreateUser(ctx context.Context, username, email, passwordHash string) (string, error) {
	var id string
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	err := s.DB.QueryRowContext(ctx, query, username, email, passwordHash).Scan(&id)
	return id, err
}

// GetUserByEmail retrieves a user by email
func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, is_verified, created_at FROM users WHERE email = $1`
	err := s.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsVerified, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &user, nil
}

// VerifyUser marks a user as verified
func (s *Store) VerifyUser(ctx context.Context, userID string) error {
	query := `UPDATE users SET is_verified = TRUE WHERE id = $1`
	_, err := s.DB.ExecContext(ctx, query, userID)
	return err
}

// CreateOTP inserts a new OTP for a user
func (s *Store) CreateOTP(ctx context.Context, userID, code string, expiresAt string) error {
	query := `INSERT INTO otps (user_id, code, expires_at) VALUES ($1, $2, $3)`
	_, err := s.DB.ExecContext(ctx, query, userID, code, expiresAt)
	return err
}

// GetLatestOTP gets the most recent OTP for a user
func (s *Store) GetLatestOTP(ctx context.Context, userID string) (*models.OTP, error) {
	var otp models.OTP
	query := `SELECT id, user_id, code, expires_at, created_at FROM otps WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`
	err := s.DB.QueryRowContext(ctx, query, userID).Scan(&otp.ID, &otp.UserID, &otp.Code, &otp.ExpiresAt, &otp.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}
