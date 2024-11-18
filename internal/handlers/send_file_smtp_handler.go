package handlers

import (
	"doodocs-archive/pkg/errors"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (h *HandlerV1) SendFileSMTPHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200 << 20) // 200 MB max memory
	if err != nil {
		errors.Error(w, err, http.StatusBadRequest, "Unable to parse form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		errors.Error(w, err, http.StatusBadRequest, "Error retrieving the file")
		return
	}
	defer file.Close()

	allowedTypes := map[string]bool{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/pdf": true,
	}
	if !allowedTypes[header.Header.Get("Content-Type")] {
		errors.Error(w, err, http.StatusBadRequest, "File type not allowed")
		return
	}

	emailsStr := r.FormValue("emails")
	if emailsStr == "" {
		errors.Error(w, err, http.StatusBadRequest, "No emails provided")
		return
	}
	emails := strings.Split(strings.ReplaceAll(emailsStr, " ", ""), ","))

	err = h.Services.SmtpService.SendToEmailList(header.Filename, file, emails)
	if err != nil {
		errors.Error(w, err, http.StatusInternalServerError, fmt.Sprintf("Error sending emails: %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "File sent successfully"})
}
