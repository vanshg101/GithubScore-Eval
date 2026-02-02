package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	GithubClientID       string
	GithubClientSecret   string
	GithubRedirectURL    string
	JWTSecret            string
	GCPProjectID         string
	FirestoreCredentials string
	MLServiceURL         string
	Environment          string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment")
	}

	cfg := &Config{
		Port:                 getEnv("PORT", "8080"),
		GithubClientID:       getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret:   getEnv("GITHUB_CLIENT_SECRET", ""),
		GithubRedirectURL:    getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/auth/github/callback"),
		JWTSecret:            getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		GCPProjectID:         getEnv("GCP_PROJECT_ID", ""),
		FirestoreCredentials: getEnv("FIRESTORE_CREDENTIALS", ""),
		MLServiceURL:         getEnv("ML_SERVICE_URL", "http://localhost:8000"),
		Environment:          getEnv("ENVIRONMENT", "development"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
