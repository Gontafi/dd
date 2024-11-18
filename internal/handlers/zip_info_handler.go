package handlers

import (
	"doodocs-archive/pkg/errors"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func (h *HandlerV1) GetZipInfoHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		errors.Error(w, err, http.StatusBadRequest, "Failed to retrieve file")
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", "uploaded-*.zip")
	if err != nil {
		errors.Error(w, err, http.StatusInternalServerError, "Failed to create temporary file")
		return
	}
	defer os.Remove(tempFile.Name())

	size, err := io.Copy(tempFile, file)
	if err != nil {
		errors.Error(w, err, http.StatusInternalServerError, "Failed to save uploaded file")
		return
	}

	result, err := h.Services.ZipService.GetZipInfo(tempFile, size)
	if err != nil {
		errors.Error(w, err, http.StatusInternalServerError, "Failed to analyze zip archive")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
