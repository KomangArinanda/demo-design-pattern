package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                     string
	DatabaseSimulatedLatency time.Duration
}

func Load() Config {
	return Config{
		Port:                     envOrDefault("APP_PORT", "7082"),
		DatabaseSimulatedLatency: time.Duration(envIntOrDefault("DB_SIMULATED_LATENCY_MS", 50)) * time.Millisecond,
	}
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envIntOrDefault(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
