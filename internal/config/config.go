package config


import (
	"log/slog"
	"os"
	"strconv"
)

// create a custom configuration type
type Config struct {
	ProxyAddr   string
	APIAddr     string
	DatabaseURL string
	TargetURL   string
	APIKey      string
	MaxTokens   int
}

// create the get env function
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// define a function to load reads environment variables and return a validated config.
// TERMINATE THE PROGRAM IF REQUIRED SECURITY KEYS ARE MISSING
func Load() *Config {
	// read the env configurations and load them into the custom type.
	maxTokensStr := getEnv("MAX_TOKENS", "4096")
	maxTokens, err := strconv.Atoi(maxTokensStr) // atoi(ascii to int) convert it to int so we use it in max tokens
	if err != nil {
		// handle error
		 slog.Warn("Invalid MAX_TOKENS value, using fallback default", "input", maxTokensStr, "fallback", 4096, "error", err)
		 maxTokens = 4096
	}
	cfg := &Config{
		ProxyAddr:   getEnv("PROXY_ADDR", ":8080"),
		APIAddr:     getEnv("API_ADDR", ":8081"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable"),
		TargetURL:   getEnv("TARGET_URL", "https://api.openai.com"),
		APIKey:      getEnv("API_KEY", ""), //Fetches the key 
		MaxTokens:   maxTokens,
	}

	// THE TERMINATION: validate required security elements before returning
	if cfg.APIKey == "" {
		// slog.Error log the issue, then exit 
		slog.Error("Critical configuration error: API_KEY is required but was not provided")
		os.Exit(1)
	}

	return cfg
}


