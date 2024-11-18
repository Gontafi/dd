package smtp_service

import (
	"bytes"
	"doodocs-archive/config"
	"doodocs-archive/pkg/utils"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/smtp"
	"sync"
)

type SMTPServiceInterface interface {
	SendToEmailList(filename string, attachment multipart.File, emails []string) error
	sendEmail(to string, attachment multipart.File, filename string) error
}

type SMTPService struct {
	cfg *config.Config
}

func SMTPInit(cfg *config.Config) SMTPServiceInterface {
	return &SMTPService{
		cfg: cfg,
	}
}

func (s *SMTPService) SendToEmailList(filename string, attachment multipart.File, emails []string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(emails))

	for _, email := range emails {
		wg.Add(1)
		email := email

		go func() {
			defer wg.Done()

			if seeker, ok := attachment.(io.Seeker); ok {
				_, err := seeker.Seek(0, 0)
				if err != nil {
					errChan <- fmt.Errorf("failed to seek attachment for %s: %v", email, err)
					return
				}
			}

			if err := s.sendEmail(email, attachment, filename); err != nil {
				errChan <- fmt.Errorf("failed to send email to %s: %v", email, err)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errors []error
	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to send emails: %v", errors)
	}

	return nil
}

func (s *SMTPService) sendEmail(to string, attachment multipart.File, filename string) error {
	auth := smtp.PlainAuth("", s.cfg.FromEmail, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	headers := map[string]string{
		"From":         s.cfg.FromEmail,
		"To":           to,
		"Subject":      "File Attachment",
		"MIME-Version": "1.0",
		"Content-Type": "multipart/mixed; boundary=boundary1",
	}

	var bodyBuffer bytes.Buffer
	bodyBuffer.WriteString("--boundary1\r\n")

	fileBuffer := &bytes.Buffer{}
	_, err := io.Copy(fileBuffer, attachment)
	if err != nil {
		return fmt.Errorf("failed to read attachment: %v", err)
	}

	mimeType := utils.DetectMimeType(fileBuffer.Bytes())

	bodyBuffer.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", mimeType, filename))
	bodyBuffer.WriteString("Content-Transfer-Encoding: base64\r\n")
	bodyBuffer.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", filename))

	encodedFile := base64.StdEncoding.EncodeToString(fileBuffer.Bytes())
	bodyBuffer.WriteString(encodedFile)
	bodyBuffer.WriteString("\r\n--boundary1--\r\n")

	emailMessage := ""
	for key, value := range headers {
		emailMessage += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	emailMessage += "\r\n" + bodyBuffer.String()

	err = smtp.SendMail(s.cfg.SMTPAddr, auth, s.cfg.FromEmail, []string{to}, []byte(emailMessage))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
