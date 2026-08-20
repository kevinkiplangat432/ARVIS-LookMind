package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/kevinkiplangat432/arvis/internal/config"
)

var ErrUnknownModel = errors.New("request model does not match any configured provider")

type bodyRequest struct {
	Model string `json:"model"`
}

// resolveProvider reads the request body far enough to find the model
// field, then returns the matched provider, the model string itself
// (needed for logging, not just routing), and an untouched copy of the
// body for forwarding.
func resolveProvider(cfg *config.Config, r *http.Request) (config.Provider, string, io.ReadCloser, error) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return config.Provider{}, "", nil, fmt.Errorf("failed to read request body: %w", err)
	}
	r.Body.Close()

	var parsed bodyRequest
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return config.Provider{}, "", nil, fmt.Errorf("failed to parse request body: %w", err)
	}
	if parsed.Model == "" {
		return config.Provider{}, "", nil, fmt.Errorf(`request body is missing a "model" field`)
	}

	provider, ok := cfg.ModelRoutes[parsed.Model]
	if !ok {
		return config.Provider{}, "", nil, ErrUnknownModel
	}

	return provider, parsed.Model, io.NopCloser(bytes.NewReader(raw)), nil
}