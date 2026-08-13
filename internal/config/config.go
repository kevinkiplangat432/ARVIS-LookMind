package config

import (
	"fmt"
	"log/slog"
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

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Load reads environment variables into a Config exactly once. It does
// not decide what's "required" — that's each command's job, since
// migrate and test don't need everything server does.
func Load() (*Config, error) {
	maxTokensStr := getEnv("MAX_TOKENS", "4096")
	maxTokens, err := strconv.Atoi(maxTokensStr)
	if err != nil {
		slog.Warn("invalid MAX_TOKENS value, using fallback default",
			"input", maxTokensStr, "fallback", 4096, "error", err)
		maxTokens = 4096
	}

	cfg := &Config{
		ProxyAddr:   getEnv("PROXY_ADDR", ":8080"),
		APIAddr:     getEnv("API_ADDR", ":8081"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable"),
		TargetURL:   getEnv("TARGET_URL", "https://api.openai.com"),
		APIKey:      getEnv("API_KEY", ""),
		MaxTokens:   maxTokens,
	}

	return cfg, nil
}

// RequireAPIKey is called only by commands that actually need to reach
// an upstream LLM provider. Right now that's just server.
func (c *Config) RequireAPIKey() error {
	if c.APIKey == "" {
		return fmt.Errorf("API_KEY is required to start the server but was not provided")
	}
	return nil
}

// String satisfies fmt.Stringer so logging a *Config anywhere (a %v, a
// slog call, a panic dump) never leaks the raw key. This is the SOC 2
// piece you remembered — it wasn't in the old code, it's new here.
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{ProxyAddr:%s APIAddr:%s DatabaseURL:%s TargetURL:%s APIKey:%s MaxTokens:%d}",
		c.ProxyAddr, c.APIAddr, c.DatabaseURL, c.TargetURL, redact(c.APIKey), c.MaxTokens,
	)
}

func redact(s string) string {
	if s == "" {
		return "(empty)"
	}
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}