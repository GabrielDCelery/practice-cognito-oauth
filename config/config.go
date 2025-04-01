package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AWSRegion           string
	CognitoClientID     string
	CognitoClientSecret string
	CognitoIssuerURL    string
	RedirectURL         string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return &Config{
		AWSRegion:           os.Getenv("AWS_REGION"),
		CognitoClientID:     os.Getenv("COGNITO_CLIENT_ID"),
		CognitoClientSecret: os.Getenv("COGNITO_CLIENT_SECRET"),
		CognitoIssuerURL:    os.Getenv("COGNITO_ISSUER_URL"),
		RedirectURL:         os.Getenv("REDIRECT_URL"),
	}
}
