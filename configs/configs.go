package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddress  string
	RedisPassword string
	RedisDB       int

	ClientBaseURL string
	ClientTimeout time.Duration
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func Load() (*Config, error) {
	// Load .env if present (ignore error in production)
	_ = godotenv.Load()

	db, err := strconv.Atoi(getEnv("DATABASE", "0"))
	if err != nil {
		return nil, err
	}

	timeout, err := time.ParseDuration(getEnv("CLIENT_TIMEOUT", "5s"))
	if err != nil {
		return nil, err
	}

	return &Config{
		RedisAddress:  getEnv("ADDRESS", "localhost:6379"),
		RedisPassword: getEnv("PASSWORD", ""),
		RedisDB:       db,
		ClientBaseURL: getEnv("CLIENT_BASE_URL", "https://api.spacexdata.com/v4"),
		ClientTimeout: timeout,
	}, nil
}