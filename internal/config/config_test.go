package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear the environment keys we are testing to ensure we get defaults
	os.Unsetenv("PROXY_ADDR")
	os.Unsetenv("MAX_TOKENS")
	
	// Set the required API key so Load() doesn't call os.Exit(1)
	os.Setenv("API_KEY", "test-key-123")
	defer os.Unsetenv("API_KEY") // Clean up after the test finishes

	cfg := Load()

	// Verify that fallbacks are applied correctly
	if cfg.ProxyAddr != ":8080" {
		t.Errorf("Expected default ProxyAddr ':8080', got '%s'", cfg.ProxyAddr)
	}
	if cfg.MaxTokens != 4096 {
		t.Errorf("Expected default MaxTokens 4096, got %d", cfg.MaxTokens)
	}
	if cfg.APIKey != "test-key-123" {
		t.Errorf("Expected APIKey 'test-key-123', got '%s'", cfg.APIKey)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	// Set custom environment variables
	os.Setenv("PROXY_ADDR", ":9090")
	os.Setenv("MAX_TOKENS", "1024")
	os.Setenv("API_KEY", "secret-key")
	
	// Always clean up environment overrides so they don't break other tests
	defer func() {
		os.Unsetenv("PROXY_ADDR")
		os.Unsetenv("MAX_TOKENS")
		os.Unsetenv("API_KEY")
	}()

	cfg := Load()

	// Verify that our custom values were successfully read and parsed
	if cfg.ProxyAddr != ":9090" {
		t.Errorf("Expected custom ProxyAddr ':9090', got '%s'", cfg.ProxyAddr)
	}
	if cfg.MaxTokens != 1024 {
		t.Errorf("Expected custom MaxTokens 1024, got %d", cfg.MaxTokens)
	}
}

func TestLoad_InvalidMaxTokensFallback(t *testing.T) {
	os.Setenv("API_KEY", "valid-key")
	os.Setenv("MAX_TOKENS", "not-a-number") // Invalid integer string
	defer func() {
		os.Unsetenv("API_KEY")
		os.Unsetenv("MAX_TOKENS")
	}()

	cfg := Load()

	// Verify that the error handling fallback mechanism safely kicks in
	if cfg.MaxTokens != 4096 {
		t.Errorf("Expected invalid MaxTokens to fall back to 4096, got %d", cfg.MaxTokens)
	}
}
