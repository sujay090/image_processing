package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"image_processing_backend/internal/db"
	"image_processing_backend/internal/mq"
	"image_processing_backend/internal/processor"
	"image_processing_backend/internal/storage"
	"image_processing_backend/internal/util"
)

// TransformJob represents a job message consumed from RabbitMQ
type TransformJob struct {
	TransformationID string                 `json:"transformation_id"`
	ImageID          string                 `json:"image_id"`
	OriginalS3Key    string                 `json:"original_s3_key"`
	Transformations  map[string]interface{} `json:"transformations"`
}

func main() {
	cfg := util.MustLoad()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	store := db.NewStore(database)

	storageClient, err := storage.NewStorageClient(
		context.Background(),
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Bucket,
	)
	if err != nil {
		log.Fatalf("Failed to connect to S3 storage: %v", err)
	}

	rabbit, err := mq.Connect(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbit.Close()

	if err := rabbit.DeclareQueue("image_transformations"); err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := rabbit.Consume("image_transformations")
	if err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}

	log.Println("Worker started. Waiting for transformation jobs...")

	for msg := range msgs {
		var job TransformJob
		if err := json.Unmarshal(msg.Body, &job); err != nil {
			log.Printf("Failed to parse job: %v", err)
			msg.Nack(false, false) // don't requeue malformed messages
			continue
		}

		log.Printf("Processing transformation %s for image %s", job.TransformationID, job.ImageID)

		err := processJob(context.Background(), store, storageClient, job)
		if err != nil {
			log.Printf("Job %s failed: %v", job.TransformationID, err)
			_ = store.UpdateTransformationStatus(context.Background(), job.TransformationID, "failed", "")
			msg.Nack(false, false)
			continue
		}

		msg.Ack(false)
		log.Printf("Job %s completed successfully", job.TransformationID)
	}
}

func processJob(ctx context.Context, store *db.Store, storageClient *storage.StorageClient, job TransformJob) error {
	// Update status to processing
	err := store.UpdateTransformationStatus(ctx, job.TransformationID, "processing", "")
	if err != nil {
		return fmt.Errorf("failed to update status to processing: %w", err)
	}

	// Download the original image from R2
	reader, err := storageClient.DownloadFile(ctx, job.OriginalS3Key)
	if err != nil {
		return fmt.Errorf("failed to download original image: %w", err)
	}
	defer reader.Close()

	imgBuf, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read image data: %w", err)
	}

	// Parse transformations into processor options
	opts, err := parseTransformOptions(job.Transformations)
	if err != nil {
		return fmt.Errorf("failed to parse transform options: %w", err)
	}

	// Process the image
	resultBuf, err := processor.Process(imgBuf, opts)
	if err != nil {
		return fmt.Errorf("image processing failed: %w", err)
	}

	// Determine output format
	outputFormat := "jpeg"
	if opts.Format != nil {
		outputFormat = *opts.Format
	}

	// Upload the processed image back to R2
	outputKey := fmt.Sprintf("transformed/%s.%s", job.TransformationID, outputFormat)
	contentType := "image/" + outputFormat

	err = storageClient.UploadFile(ctx, outputKey, bytes.NewReader(resultBuf), contentType)
	if err != nil {
		return fmt.Errorf("failed to upload transformed image: %w", err)
	}

	// Update transformation record as completed
	err = store.UpdateTransformationStatus(ctx, job.TransformationID, "completed", outputKey)
	if err != nil {
		return fmt.Errorf("failed to update status to completed: %w", err)
	}

	return nil
}

func parseTransformOptions(raw map[string]interface{}) (processor.TransformOptions, error) {
	var opts processor.TransformOptions

	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return opts, err
	}
	err = json.Unmarshal(jsonBytes, &opts)
	return opts, err
}
