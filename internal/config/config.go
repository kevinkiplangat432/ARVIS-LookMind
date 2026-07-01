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


// code revew 1 
// 
// Lines 11–18 (type Config struct): This is your "Source of Truth." By centralizing configuration into a struct, you prevent "magic strings" from being scattered across your codebase.
// Lines 21–26 (getEnv): A helper to provide sensible defaults. This makes the developer experience (DX) better because the app "just works" locally without a .env file
// Lines 43–47 (The TERMINATION): This is the Fail-Fast Principle. It is much better to crash at second 0 than to start the server and have 1,000 requests fail 5 minutes later because of a missing key.

// soc 2 type ii pass/fail analysis
// verdict fail 
// Secret Masking (Fail): If you ever use fmt.Printf("%+v", cfg) for debugging, your APIKey and DatabaseURL (which contains the password) will be printed in plain text to your logs.Remediation: Implement the Stringer interface to mask sensitive values (shown in the refactor below).
// Audit Trail (Pass): You are using slog.Error for the missing API key. This provides the "Evidence of Control" auditors love to see.
// Environment Isolation (Warning): Your defaults (like the localhost Postgres URL) are great for local dev, but in a SOC 2 environment, you must ensure these defaults cannot accidentally be used in Production.

//scaling to 1000 RPS 
// Data Leakage (High Risk): Your DatabaseURL default contains postgres://arvis:arvis.... If an error occurs during connection, many Go drivers return the full connection string in the error message.
// Performance (Pass): Configuration is usually loaded once at startup. This file isn't on the "hot path" of your 1,000 RPS, so its performance is fine. However, Validation is the key to stability.