package db

import (
	"context"

	"image_processing_backend/internal/models"
)

// CreateImage inserts a new image record
func (s *Store) CreateImage(ctx context.Context, img *models.Image) error {
	query := `
		INSERT INTO images (id, user_id, original_filename, s3_key, format, size, width, height, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := s.DB.ExecContext(ctx, query,
		img.ID, img.UserID, img.OriginalFilename, img.S3Key, img.Format, img.Size, img.Width, img.Height, img.Status,
	)
	return err
}

// GetImage retrieves an image by its ID and ensures it belongs to the user
func (s *Store) GetImage(ctx context.Context, id, userID string) (*models.Image, error) {
	var img models.Image
	query := `
		SELECT id, user_id, original_filename, s3_key, format, size, width, height, status, created_at
		FROM images
		WHERE id = $1 AND user_id = $2
	`
	err := s.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&img.ID, &img.UserID, &img.OriginalFilename, &img.S3Key, &img.Format,
		&img.Size, &img.Width, &img.Height, &img.Status, &img.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

// ListImages retrieves a paginated list of images for a user
func (s *Store) ListImages(ctx context.Context, userID string, limit, offset int) ([]*models.Image, error) {
	query := `
		SELECT id, user_id, original_filename, s3_key, format, size, width, height, status, created_at
		FROM images
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := s.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*models.Image
	for rows.Next() {
		var img models.Image
		err := rows.Scan(
			&img.ID, &img.UserID, &img.OriginalFilename, &img.S3Key, &img.Format,
			&img.Size, &img.Width, &img.Height, &img.Status, &img.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, &img)
	}
	return images, rows.Err()
}

// CreateTransformation creates a new transformation job
func (s *Store) CreateTransformation(ctx context.Context, tf *models.Transformation) error {
	query := `
		INSERT INTO transformations (id, image_id, transform_params, s3_key, status)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := s.DB.ExecContext(ctx, query, tf.ID, tf.ImageID, tf.TransformParams, tf.S3Key, tf.Status)
	return err
}

// GetTransformation retrieves a transformation by ID
func (s *Store) GetTransformation(ctx context.Context, id string) (*models.Transformation, error) {
	var tf models.Transformation
	query := `
		SELECT id, image_id, transform_params, s3_key, status, created_at
		FROM transformations
		WHERE id = $1
	`
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&tf.ID, &tf.ImageID, &tf.TransformParams, &tf.S3Key, &tf.Status, &tf.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &tf, nil
}
