package config

import (
	"os"
	"strconv"
	"time"
)

// HTTPConfig contains configs for HTTP
type HTTPConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Config contains all configs of multiplexer
type Config struct {
	HTTP HTTPConfig
}

// NewConfig create instance of Config
func NewConfig() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Host: getStringEnv("HTTP_HOST", "localhost"),
			Port: getIntEnv("HTTP_PORT", 8080),
		},
	}
}

// getStringEnv read an environment as string or return a default value
func getStringEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

// getIntEnv read an environment as int or return a default value
func getIntEnv(key string, defaultValue int) int {
	str := getStringEnv(key, "")

	if value, err := strconv.Atoi(str); err == nil {
		return value
	}

	return defaultValue
}

// getTimeDurationEnv read an environment as time.Duration or return a default value
func getTimeDurationEnv(key string, defaultValue time.Duration) time.Duration {
	str := getStringEnv(key, "")

	if value, err := time.ParseDuration(str); err == nil {
		return value
	}

	return defaultValue
}