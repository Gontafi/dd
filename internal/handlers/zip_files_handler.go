package handlers

import (
	"doodocs-archive/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func (h *HandlerV1) ZipFilesHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(200 << 20) // 200 MB
	if err != nil {
		errors.Error(w, err, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	allowedMimeTypes := map[string]bool{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/xml": true,
		"image/jpeg":      true,
		"image/png":       true,
	}

	var files []multipart.File
	var filenames []string

	for _, fileHeader := range r.MultipartForm.File["files[]"] {
		if !allowedMimeTypes[fileHeader.Header.Get("Content-Type")] {
			errors.Error(w, err, http.StatusBadRequest, "File type not allowed: "+fileHeader.Filename)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			errors.Error(w, err, http.StatusBadRequest, "Failed to open uploaded file: "+fileHeader.Filename)
			return
		}
		defer file.Close()

		files = append(files, file)
		filenames = append(filenames, filepath.Base(fileHeader.Filename))
	}

	if len(files) == 0 {
		errors.Error(w, err, http.StatusBadRequest, "No files were uploaded")
		return
	}

	zipReader, err := h.Services.ZipService.CreateZip(files, filenames)
	if err != nil {
		errors.Error(w, err, http.StatusInternalServerError, "Failed to create zip archive")
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=archive.zip")
	_, err = io.Copy(w, zipReader)
	if err != nil {
		errors.Error(w, err, http.StatusInternalServerError, "Failed to send zip archive")
		return
	}
}
