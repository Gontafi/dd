package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	FromEmail    string
	SMTPPassword string
	SMTPHost     string
	SMTPAddr     string
}

func LoadConfig() *Config {
	loadEnv()

	return &Config{
		FromEmail:    getEnv("FROM_EMAIL", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPHost:     getEnv("FROM_EMAIL_SMTP", ""),
		SMTPAddr:     getEnv("SMTP_ADDR", ""),
	}
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found in /src, trying /etc/secrets/.env")

		if err := godotenv.Load("/etc/secrets/.env"); err != nil {
			log.Println(".env not found in /etc/secrets either, using system environment variables")
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
