package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey                string
	GRPCPort              string
	OpticOddsAPIBaseURL   string
	OddsAIAPIBaseURL      string
	DBUsername            string
	DBPassword            string
	DBName                string
	DBHost                string
	DBPort                string
	RedisHost             string
	RedisPassword         string
	APICorsAllowedOrigins string
}

func LoadConfig() *Config {
	// Load .env file from the root directory
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		APIKey:                getEnv("API_KEY", ""),
		GRPCPort:              getEnv("GRPC_PORT", ":50051"),
		OpticOddsAPIBaseURL:   getEnv("OPTIC_ODDS_API_BASE_URL", "https://api.opticodds.com"),
		OddsAIAPIBaseURL:      getEnv("ODDS_AI_API_BASE_URL", "https://demofeed.betapi.win/FeedApi"),
		DBUsername:            getEnv("DB_USERNAME", ""),
		DBPassword:            getEnv("DB_PASSWORD", ""),
		DBName:                getEnv("DB_NAME", ""),
		DBHost:                getEnv("DB_HOST", "localhost"),
		DBPort:                getEnv("DB_PORT", "5432"),
		RedisHost:             getEnv("REDIS_HOST", "51.159.19.90:6379"),
		RedisPassword:         getEnv("REDIS_PASSWORD", "Meganet123!@#"),
		APICorsAllowedOrigins: getEnv("API_CORS_ALLOWED_ORIGINS", "http://*,https://*"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Environment variable %s not set. Defaulting to %s.", key, fallback)
	return fallback
}
