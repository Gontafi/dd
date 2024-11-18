package handlers

import (
	"doodocs-archive/config"
	"doodocs-archive/internal/services"
	"net/http"
)

type HandlerInterface interface {
	SendFileSMTPHandler(http.ResponseWriter, *http.Request)
	ZipFilesHandler(w http.ResponseWriter, r *http.Request)
	GetZipInfoHandler(w http.ResponseWriter, r *http.Request)
}

type HandlerV1 struct {
	Services *services.Services
}

func NewHandlerV1(cfg *config.Config) HandlerInterface {
	return &HandlerV1{
		Services: services.ServiceInit(cfg),
	}
}
