package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	ProxyAddr     string
	APIAddr       string
	DatabaseURL   string
	ProvidersPath string
	Providers     []Provider
	ModelRoutes   map[string]Provider
	MaxTokens     int
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Load reads environment variables into a Config exactly once. It does
// not decide what's "required" — that's each command's job, since
// migrate and test don't need everything server does. Notably, it does
// NOT load providers.yaml — that's ResolveProviders, called explicitly
// only by commands that actually route traffic.
func Load() (*Config, error) {
	maxTokensStr := getEnv("MAX_TOKENS", "4096")
	maxTokens, err := strconv.Atoi(maxTokensStr)
	if err != nil {
		slog.Warn("invalid MAX_TOKENS value, using fallback default",
			"input", maxTokensStr, "fallback", 4096, "error", err)
		maxTokens = 4096
	}

	cfg := &Config{
		ProxyAddr:     getEnv("PROXY_ADDR", ":8080"),
		APIAddr:       getEnv("API_ADDR", ":8081"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable"),
		ProvidersPath: getEnv("PROVIDERS_FILE", "providers.yaml"),
		MaxTokens:     maxTokens,
	}

	return cfg, nil
}

// ResolveProviders loads and validates the providers file, then builds
// ModelRoutes — a flat model-name-to-provider map — so Phase 5's proxy
// can route in O(1) per request instead of scanning every provider's
// model list on every call.
func (c *Config) ResolveProviders() error {
	providers, err := LoadProviders(c.ProvidersPath)
	if err != nil {
		return fmt.Errorf("providers are required to start the server: %w", err)
	}

	routes := make(map[string]Provider)
	for _, p := range providers {
		for _, model := range p.Models {
			if existing, ok := routes[model]; ok {
				return fmt.Errorf("model %q is claimed by both %q and %q — each model must map to exactly one provider", model, existing.Name, p.Name)
			}
			routes[model] = p
		}
	}

	c.Providers = providers
	c.ModelRoutes = routes
	return nil
}

// String satisfies fmt.Stringer so logging a *Config anywhere (a %v, a
// slog call, a panic dump) never leaks provider keys — only names.
func (c *Config) String() string {
	names := make([]string, len(c.Providers))
	for i, p := range c.Providers {
		names[i] = p.Name
	}
	return fmt.Sprintf(
		"Config{ProxyAddr:%s APIAddr:%s DatabaseURL:%s ProvidersPath:%s Providers:%v MaxTokens:%d}",
		c.ProxyAddr, c.APIAddr, c.DatabaseURL, c.ProvidersPath, names, c.MaxTokens,
	)
}

