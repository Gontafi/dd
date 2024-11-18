package handlers

import (
	"archive/zip"
	er "doodocs-archive/pkg/errors"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func (h *HandlerV1) GetZipInfoHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		er.Error(w, err, http.StatusBadRequest, "Failed to retrieve file")
		return
	}
	defer file.Close()

	filename := header.Filename

	tempFile, err := os.CreateTemp("", "uploaded-*.zip")
	if err != nil {
		er.Error(w, err, http.StatusInternalServerError, "Failed to create temporary file")
		return
	}
	defer os.Remove(tempFile.Name())

	size, err := io.Copy(tempFile, file)
	if err != nil {
		er.Error(w, err, http.StatusInternalServerError, "Failed to save uploaded file")
		return
	}

	// Pass the filename along with the tempFile and size to the service
	result, err := h.Services.ZipService.GetZipInfo(filename, tempFile, size)
	if err != nil {
		if errors.Is(err, zip.ErrFormat) {
			er.Error(w, err, http.StatusBadRequest, err.Error())
			return
		}

		er.Error(w, err, http.StatusInternalServerError, "Failed to analyze zip archive")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
