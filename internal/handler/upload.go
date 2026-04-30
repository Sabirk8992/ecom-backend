package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sabirk8992/ecom-backend/internal/storage"
)

type UploadHandler struct {
	Storage *storage.S3Storage
}

func NewUploadHandler(s *storage.S3Storage) *UploadHandler {
	return &UploadHandler{Storage: s}
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Max 5MB
	r.ParseMultipartForm(5 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), header.Filename)
	contentType := header.Header.Get("Content-Type")

	url, err := h.Storage.Upload(file, fileName, contentType)
	if err != nil {
		http.Error(w, "failed to upload file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}
