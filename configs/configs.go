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

	CacheTTL time.Duration
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func Load() (*Config, error) {
	_ = godotenv.Load() // safe for local, ignored in container

	db, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, err
	}

	timeout, err := time.ParseDuration(getEnv("CLIENT_TIMEOUT", "5s"))
	if err != nil {
		return nil, err
	}

	ttl, err := time.ParseDuration(getEnv("CACHE_TTL", "60s"))

	return &Config{
		RedisAddress:  getEnv("REDIS_ADDRESS", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       db,
		ClientBaseURL: getEnv("CLIENT_BASE_URL", "https://api.spacexdata.com/v4"),
		ClientTimeout: timeout,
		CacheTTL: ttl,
	}, nil
}