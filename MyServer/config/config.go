package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configurations
type Config struct {
	DatabaseURL         string
	JWTSecret           string
	ServerPort          string
	EncryptionNonceSize int
}

// LoadConfig loads configuration from environment variables or .env file
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, loading from environment variables.")
	}

	cfg := &Config{
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/magicgate?sslmode=disable"),
		JWTSecret:           getEnv("JWT_SECRET", "supersecretjwtkey"), // IMPORTANT: Change this in production!
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		EncryptionNonceSize: getEnvAsInt("ENCRYPTION_NONCE_SIZE", 12), // GCM recommended nonce size
	}

	if cfg.JWTSecret == "supersecretjwtkey" {
		log.Println("WARNING: Using default JWT_SECRET. Please set a strong secret in production.")
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
