package handler

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"image_processing_backend/internal/db"
	"image_processing_backend/internal/models"
	"image_processing_backend/internal/mq"
	"image_processing_backend/internal/storage"
	"image_processing_backend/internal/middleware"
)

type ImageHandler struct {
	store   *db.Store
	storage *storage.StorageClient
	rabbit  *mq.RabbitMQ
}

func NewImageHandler(store *db.Store, storage *storage.StorageClient, rabbit *mq.RabbitMQ) *ImageHandler {
	return &ImageHandler{
		store:   store,
		storage: storage,
		rabbit:  rabbit,
	}
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		sendError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB max memory
	if err != nil {
		sendError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		sendError(w, http.StatusBadRequest, "Image file is required")
		return
	}
	defer file.Close()

	// Peek at the image configuration to get width/height and format
	imgConfig, format, err := image.DecodeConfig(file)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid image format")
		return
	}

	// Reset file pointer to beginning for upload
	file.Seek(0, 0)

	fileID := uuid.New().String()
	s3Key := fmt.Sprintf("%s/%s.%s", userID, fileID, format)

	contentType := "image/" + format

	err = h.storage.UploadFile(r.Context(), s3Key, file, contentType)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to upload image to storage")
		return
	}

	imgRecord := &models.Image{
		ID:               fileID,
		UserID:           userID,
		OriginalFilename: header.Filename,
		S3Key:            s3Key,
		Format:           format,
		Size:             header.Size,
		Width:            imgConfig.Width,
		Height:           imgConfig.Height,
		Status:           "ready",
	}

	err = h.store.CreateImage(r.Context(), imgRecord)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to save image record")
		return
	}

	sendJSON(w, http.StatusCreated, imgRecord)
}

func (h *ImageHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		sendError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := (page - 1) * limit

	images, err := h.store.ListImages(r.Context(), userID, limit, offset)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch images")
		return
	}

	sendJSON(w, http.StatusOK, images)
}

func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		sendError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	// In Go 1.22+ standard mux, path variables can be accessed via r.PathValue
	imageID := r.PathValue("id")
	if imageID == "" {
		sendError(w, http.StatusBadRequest, "Image ID is required")
		return
	}

	img, err := h.store.GetImage(r.Context(), imageID, userID)
	if err != nil {
		sendError(w, http.StatusNotFound, "Image not found")
		return
	}

	sendJSON(w, http.StatusOK, img)
}

type TransformRequest struct {
	Transformations map[string]interface{} `json:"transformations"`
}

func (h *ImageHandler) TransformImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		sendError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	imageID := r.PathValue("id")
	if imageID == "" {
		sendError(w, http.StatusBadRequest, "Image ID is required")
		return
	}

	var req TransformRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	img, err := h.store.GetImage(r.Context(), imageID, userID)
	if err != nil {
		sendError(w, http.StatusNotFound, "Image not found")
		return
	}

	paramsBytes, err := json.Marshal(req.Transformations)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to parse transformations")
		return
	}

	tfID := uuid.New().String()
	tf := &models.Transformation{
		ID:              tfID,
		ImageID:         img.ID,
		TransformParams: string(paramsBytes),
		S3Key:           "",
		Status:          "pending",
	}

	err = h.store.CreateTransformation(r.Context(), tf)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create transformation record")
		return
	}

	// Publish to RabbitMQ
	jobPayload := map[string]interface{}{
		"transformation_id": tf.ID,
		"image_id":          img.ID,
		"original_s3_key":   img.S3Key,
		"transformations":   req.Transformations,
	}
	jobBytes, _ := json.Marshal(jobPayload)

	err = h.rabbit.Publish(r.Context(), "image_transformations", jobBytes)
	if err != nil {
		// ideally we should mark tf as failed here in DB, but for simplicity:
		sendError(w, http.StatusInternalServerError, "Failed to queue transformation job")
		return
	}

	sendJSON(w, http.StatusAccepted, map[string]string{
		"message":           "Transformation job queued",
		"transformation_id": tf.ID,
	})
}

func (h *ImageHandler) GetTransformationStatus(w http.ResponseWriter, r *http.Request) {
	tfID := r.PathValue("id")
	if tfID == "" {
		sendError(w, http.StatusBadRequest, "Transformation ID is required")
		return
	}

	tf, err := h.store.GetTransformation(r.Context(), tfID)
	if err != nil {
		sendError(w, http.StatusNotFound, "Transformation not found")
		return
	}

	sendJSON(w, http.StatusOK, tf)
}
