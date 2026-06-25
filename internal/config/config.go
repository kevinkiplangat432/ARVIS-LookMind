package config


import (
	"os"
	"strconv"
)

type Config struct {
	ProxyAddr   string
	APIAddr     string
	DatabaseURL string
	TargetURL   string
	APIKey      string
	MaxTokens   int
}

func Load() *Config {
	maxTokens, _ := strconv.Atoi(getEnv("MAX_TOKENS", "4096"))
	return &Config{
		ProxyAddr:   getEnv("PROXY_ADDR", ":8080"),
		APIAddr:     getEnv("API_ADDR", ":8081"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable"),
		TargetURL:   getEnv("TARGET_URL", "https://api.openai.com"),
		APIKey:      getEnv("API_KEY", ""),
		MaxTokens:   maxTokens,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
