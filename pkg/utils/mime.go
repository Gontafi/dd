package utils

import (
	"archive/zip"
	"net/http"
)

func DetectMimeTypeZip(file *zip.File) string {

	if file == nil {
		return "application/octet-stream"
	}

	rc, err := file.Open()
	if err != nil {
		return "application/octet-stream"
	}
	defer rc.Close()

	buffer := make([]byte, 512)
	n, err := rc.Read(buffer)
	if err != nil && n == 0 {
		return "application/octet-stream"
	}

	return http.DetectContentType(buffer[:n])
}

func DetectMimeType(data []byte) string {

	if data == nil || len(data) == 0 {
		return "application/octet-stream"
	}
	return http.DetectContentType(data)
}
