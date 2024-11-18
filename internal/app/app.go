package app

import (
	"doodocs-archive/config"
	"doodocs-archive/internal/handlers"
	"doodocs-archive/internal/middleware"
	"log"
	"net/http"
)

func RunApp() {
	cfg := config.LoadConfig()

	v1 := handlers.NewHandlerV1(cfg)

	http.HandleFunc("POST /api/archive/information", middleware.LoggerMiddleware(v1.GetZipInfoHandler))
	http.HandleFunc("POST /api/archive/files", middleware.LoggerMiddleware(v1.ZipFilesHandler))
	http.HandleFunc("POST /api/mail/file", middleware.LoggerMiddleware(v1.SendFileSMTPHandler))

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
