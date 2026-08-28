package main

import (
	"context"
	"log"
	"net/http"

	"image_processing_backend/internal/db"
	"image_processing_backend/internal/mq"
	"image_processing_backend/internal/router"
	"image_processing_backend/internal/storage"
	"image_processing_backend/internal/util"
)

func main() {
	cfg := util.MustLoad()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	store := db.NewStore(database)

	// Initialize Cloudflare R2 (S3) storage client
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

	// Initialize RabbitMQ connection
	rabbit, err := mq.Connect(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbit.Close()

	// Declare the transformation queue
	if err := rabbit.DeclareQueue("image_transformations"); err != nil {
		log.Fatalf("Failed to declare RabbitMQ queue: %v", err)
	}

	mux := router.SetupRouter(store, cfg, storageClient, rabbit)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
