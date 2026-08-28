package router

import (
	"net/http"

	"image_processing_backend/internal/db"
	"image_processing_backend/internal/handler"
	"image_processing_backend/internal/middleware"
	"image_processing_backend/internal/mq"
	"image_processing_backend/internal/storage"
	"image_processing_backend/internal/util"
)

func SetupRouter(store *db.Store, cfg *util.Config, storageClient *storage.StorageClient, rabbit *mq.RabbitMQ) *http.ServeMux {
	mux := http.NewServeMux()
	authHandler := handler.NewAuthHandler(store, cfg.JWTSecret)
	imageHandler := handler.NewImageHandler(store, storageClient, rabbit)

	// Auth routes
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/verify-email", authHandler.VerifyEmail)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	// Protected image routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("POST /api/v1/images", imageHandler.UploadImage)
	protectedMux.HandleFunc("GET /api/v1/images", imageHandler.ListImages)
	protectedMux.HandleFunc("GET /api/v1/images/{id}", imageHandler.GetImage)
	protectedMux.HandleFunc("POST /api/v1/images/{id}/transform", imageHandler.TransformImage)
	protectedMux.HandleFunc("GET /api/v1/transformations/{id}", imageHandler.GetTransformationStatus)

	mux.Handle("/api/v1/images/", middleware.AuthMiddleware(cfg.JWTSecret, protectedMux))
	mux.Handle("/api/v1/images", middleware.AuthMiddleware(cfg.JWTSecret, protectedMux))
	mux.Handle("/api/v1/transformations/", middleware.AuthMiddleware(cfg.JWTSecret, protectedMux))

	return mux
}
