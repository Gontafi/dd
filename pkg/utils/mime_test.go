package utils

import (
	"archive/zip"
	"bytes"
	"net/http"
	"testing"
)

func TestDetectMimeTypeZip(t *testing.T) {
	// Create a test ZIP file in memory
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Add a text file to the ZIP
	f, err := w.Create("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write([]byte("This is a test file."))
	if err != nil {
		t.Fatal(err)
	}

	// Close the ZIP writer
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Open the ZIP file
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}

	// Test DetectMimeTypeZip
	mimeType := DetectMimeTypeZip(r.File[0])
	if mimeType != "text/plain; charset=utf-8" {
		t.Errorf("Expected mime type 'text/plain; charset=utf-8', got '%s'", mimeType)
	}
}

func TestDetectMimeType(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "PNG image",
			data:     []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			expected: "image/png",
		},
		{
			name:     "JPEG image",
			data:     []byte{0xFF, 0xD8, 0xFF},
			expected: "image/jpeg",
		},
		{
			name:     "Plain text",
			data:     []byte("This is a plain text file."),
			expected: "text/plain; charset=utf-8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mimeType := DetectMimeType(tc.data)
			if mimeType != tc.expected {
				t.Errorf("Expected mime type '%s', got '%s'", tc.expected, mimeType)
			}
		})
	}
}

func TestDetectMimeTypeZipEmptyFile(t *testing.T) {
	// Create a test ZIP file in memory with an empty file
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Add an empty file to the ZIP
	_, err := w.Create("empty.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Close the ZIP writer
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Open the ZIP file
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}

	// Test DetectMimeTypeZip with empty file
	mimeType := DetectMimeTypeZip(r.File[0])
	if mimeType != "application/octet-stream" {
		t.Errorf("Expected 'application/octet-stream' for empty file, got '%s'", mimeType)
	}
}

func TestDetectMimeTypeLargeFile(t *testing.T) {
	// Create a large byte slice (larger than 512 bytes)
	largeData := make([]byte, 1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// First 512 bytes should be used for detection
	expected := http.DetectContentType(largeData[:512])

	mimeType := DetectMimeType(largeData)
	if mimeType != expected {
		t.Errorf("Expected mime type '%s', got '%s'", expected, mimeType)
	}
}

func TestDetectMimeTypeEdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "Empty data",
			data:     []byte{},
			expected: "application/octet-stream",
		},
		{
			name:     "Very short data",
			data:     []byte{0x00, 0x01},
			expected: "application/octet-stream",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mimeType := DetectMimeType(tc.data)
			if mimeType != tc.expected {
				t.Errorf("Expected mime type '%s', got '%s'", tc.expected, mimeType)
			}
		})
	}
}

func TestDetectMimeTypeZipLargeFile(t *testing.T) {
	// Create a test ZIP file in memory with a large file
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Add a large file to the ZIP
	f, err := w.Create("large.bin")
	if err != nil {
		t.Fatal(err)
	}

	// Write 1MB of random data
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}
	_, err = f.Write(largeData)
	if err != nil {
		t.Fatal(err)
	}

	// Close the ZIP writer
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Open the ZIP file
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}

	// Test DetectMimeTypeZip with large file
	mimeType := DetectMimeTypeZip(r.File[0])
	expected := http.DetectContentType(largeData[:512])
	if mimeType != expected {
		t.Errorf("Expected mime type '%s', got '%s'", expected, mimeType)
	}
}
