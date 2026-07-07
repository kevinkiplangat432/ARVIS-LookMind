package config

import (
	"fmt"
	"net"
	"log/slog"
	"net/url"
	"os"
	"strconv"
)

// Config contains all application configurations
// once loaded allconfiguration should be treated as immutable

type Config struct {
	ProxyAddr   string
	APIAddr     string
	DatabaseURL string
	TargetURL   string
	APIKey      string
	MaxTokens   int
}

// getEnv returns the environments variable if it exist,
// otherwise it return the fallback value 
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
// define load :  to load, validate and return application configuration

// define a function to load reads environment variables and return a validated config.
// TERMINATE THE PROGRAM IF REQUIRED SECURITY KEYS ARE MISSING
func Load() (*Config, error) {
	// read the env configurations and load them into the custom type.
	maxTokensStr := getEnv("MAX_TOKENS", "4096")
	maxTokens, err := strconv.Atoi(maxTokensStr) // atoi(ascii to int) convert it to int so we use it in max tokens
	if err != nil {
		// handle error
		 slog.Warn("Invalid MAX_TOKENS value. using fallback default.",
		  "input", maxTokensStr, 
		  "fallback", 4096, 
		  "error", err,
		)
		 maxTokens = 4096
	}
	cfg := &Config{
		ProxyAddr:   getEnv("PROXY_ADDR", ":8080"),
		APIAddr:     getEnv("API_ADDR", ":8081"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		TargetURL:   getEnv("TARGET_URL", "https://api.openai.com"),
		APIKey:      getEnv("API_KEY", ""), //Fetches the key 
		MaxTokens:   maxTokens,
	}

	// THE TERMINATION: validate required security elements before returning
	if err := validate(cfg); err != nil {
		return nil, err
	}
	// if cfg.APIKey == "" {
	// 	// slog.Error log the issue, then exit 
	// 	slog.Error("Critical configuration error: API_KEY is required but was not provided")
	// 	os.Exit(1)
	// }

	return cfg, nil
}

// validate performs all configuration validation
func validate(cfg *Config) error{
	if cfg.APIKey == "" {
		return fmt.Errorf("API_key is required")
	}

	if cfg.DatabaseURL == "" {
		return fmt.Errorf("Database_URL is required")
	}

	if _, err := url.ParseRequestURI(cfg.TargetURL); err != nil {
		return fmt.Errorf("invalid Target_URL: %w", err)
	}
	
	
	if _, err := url.Parse(cfg.DatabaseURL); err != nil {
		return fmt.Errorf("invalid DATABASE_URL: %w", err)
	}

	if _, err := net.ResolveTCPAddr("tcp", cfg.ProxyAddr); err != nil {
		return fmt.Errorf("invalid PROXY_ADDR: %w", err)
	}

	if _, err := net.ResolveTCPAddr("tcp", cfg.APIAddr); err != nil {
		return fmt.Errorf("invalid API_ADDR: %w", err)
	}

	if cfg.MaxTokens <= 0 || cfg.MaxTokens > 200000 {
		return fmt.Errorf("MAX_TOKENS must be between 1 and 200000")
	}

	return nil
}


// string mask sensitive fields. to prevent them from appearing in logs
func (c Config) String() string{
	apiKey := "********"
	if c.APIKey == ""{
		apiKey = "<empty>"
	}

	db := "<hidden>"

	if c.DatabaseURL ==""{
		db = "<empty>"
	}
	return fmt.Sprintf(
		"Config{ProxyAddr:%q APIAddr:%q DatabaseURL:%q TargetURL:%q APIKey:%q MaxTokens:%d}",
		c.ProxyAddr,
		c.APIAddr,
		db,
		c.TargetURL,
		apiKey,
		c.MaxTokens,
	)
}


// code revew 1 
// soc 2 type ii pass/fail analysis
// verdict fail 
// Secret Masking (Fail): If you ever use fmt.Printf("%+v", cfg) for debugging, your APIKey and DatabaseURL (which contains the password) will be printed in plain text to your logs.Remediation: Implement the Stringer interface to mask sensitive values 

// Audit Trail (Pass): You are using slog.Error for the missing API key. This provides the "Evidence of Control" auditors love to see.

// Environment Isolation: Your defaults (like the localhost Postgres URL) are great for local dev, but in a SOC 2 environment, you must ensure these defaults cannot accidentally be used in Production.

//scaling to 10000 RPS 
// Data Leakage (High Risk): Your DatabaseURL default contains postgres://arvis:arvis.... If an error occurs during connection, many Go drivers return the full connection string in the error message.

// Performance (Pass): Configuration is usually loaded once at startup. This file isn't on the "hot path" of your 1,000 RPS, so its performance is fine. However, Validation is the key to stability.


//soc 2 type ii review

/* cc6 - logical access 
api key validation Pass 


*/

/*
CC7 -- Monitoring
Structured Logging
pass
*/


/*
CC8 -- change management 
configuration centralized 
pass
*/

/*
secret Protection 
fail
reason: sensitive values may be logged accidentally
*/


/*
Production isolation
needs improvement
reason : Development defaults should not be usable in prod

*/

/*
configuration validation
Partial pass:
Only API_KEY is validated
Other Important values are not
--- malformed database URLS
	invalid HTTP adresses
	invalid TargetURL
	negative MaxTokens

*/