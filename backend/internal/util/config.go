package util

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	RabbitMQURL string
	S3Bucket    string
	S3AccessKey string
	S3SecretKey string
	S3Endpoint  string
}

// MustLoad loads environment variables and returns a Config struct.
// It fatally logs (graceful shutdown of init) if required variables are missing.
func MustLoad() *Config {
	cfg := &Config{
		Port:        getEnvOrDefault("PORT", "8080"),
		DatabaseURL: getEnvOrFatal("DATABASE_URL"),
		JWTSecret:   getEnvOrFatal("JWT_SECRET"),
		RabbitMQURL: getEnvOrDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		S3Bucket:    getEnvOrFatal("S3_BUCKET"),
		S3AccessKey: getEnvOrFatal("S3_ACCESS_KEY"),
		S3SecretKey: getEnvOrFatal("S3_SECRET_KEY"),
		S3Endpoint:  getEnvOrFatal("S3_ENDPOINT"),
	}
	return cfg
}

func getEnvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Fatal: missing required environment variable: %s", key)
	}
	return value
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
