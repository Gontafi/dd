package services

import (
	"doodocs-archive/config"
	"doodocs-archive/internal/services/archive_service"
	"doodocs-archive/internal/services/smtp_service"
)

type Services struct {
	ZipService  archive_service.ZipServiceInterface
	SmtpService smtp_service.SMTPServiceInterface
}

func ServiceInit(cfg *config.Config) *Services {
	return &Services{
		ZipService:  archive_service.ZipInit(),
		SmtpService: smtp_service.SMTPInit(cfg),
	}
}
