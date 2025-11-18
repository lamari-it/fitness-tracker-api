package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSLMode           string
	JWTSecret           string
	JWTExpires          string
	GoogleClientID      string
	GoogleClientSecret  string
	GoogleRedirectURL   string
	AppleClientID       string
	AppleTeamID         string
	AppleKeyID          string
	ApplePrivateKeyPath string
	AppleRedirectURL    string
	UseMigrations       bool
	Environment         string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", ""),
		DBName:              getEnv("DB_NAME", "fitflow"),
		DBSSLMode:           getEnv("DB_SSLMODE", "disable"),
		JWTSecret:           getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpires:          getEnv("JWT_EXPIRES_IN", "24h"),
		GoogleClientID:      getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:  getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:   getEnv("GOOGLE_REDIRECT_URL", ""),
		AppleClientID:       getEnv("APPLE_CLIENT_ID", ""),
		AppleTeamID:         getEnv("APPLE_TEAM_ID", ""),
		AppleKeyID:          getEnv("APPLE_KEY_ID", ""),
		ApplePrivateKeyPath: getEnv("APPLE_PRIVATE_KEY_PATH", ""),
		AppleRedirectURL:    getEnv("APPLE_REDIRECT_URL", ""),
		UseMigrations:       getEnv("USE_MIGRATIONS", "false") == "true",
		Environment:         getEnv("APP_ENV", "development"),
	}
}

func Load() (*Config, error) {
	LoadConfig()
	return AppConfig, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
