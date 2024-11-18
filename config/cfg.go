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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		FromEmail:    getEnv("FROM_EMAIL", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPHost:     getEnv("FROM_EMAIL_SMTP", ""),
		SMTPAddr:     getEnv("SMTP_ADDR", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
