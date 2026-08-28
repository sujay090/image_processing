package db

import (
	"context"
)

// UpdateTransformationStatus updates the status and s3_key of a transformation
func (s *Store) UpdateTransformationStatus(ctx context.Context, id, status, s3Key string) error {
	query := `UPDATE transformations SET status = $1, s3_key = $2 WHERE id = $3`
	_, err := s.DB.ExecContext(ctx, query, status, s3Key, id)
	return err
}

// UpdateImageStatus updates the status of an image
func (s *Store) UpdateImageStatus(ctx context.Context, id, status string) error {
	query := `UPDATE images SET status = $1 WHERE id = $2`
	_, err := s.DB.ExecContext(ctx, query, status, id)
	return err
}
