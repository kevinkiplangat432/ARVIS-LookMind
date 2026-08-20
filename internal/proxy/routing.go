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

// bodyRequest is deliberately minimal — ARVIS only needs the model
// field to route, it doesn't need to understand a provider's full
// request schema.
type bodyRequest struct {
	Model string `json:"model"`
}

// resolveProvider reads the request body far enough to find the model
// field, then returns both the matched provider and an untouched copy
// of the body for forwarding. An http.Request's body can only be read
// once, so this hands back a fresh io.ReadCloser built from the bytes
// already consumed.
func resolveProvider(cfg *config.Config, r *http.Request) (config.Provider, io.ReadCloser, error) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return config.Provider{}, nil, fmt.Errorf("failed to read request body: %w", err)
	}
	r.Body.Close()

	var parsed bodyRequest
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return config.Provider{}, nil, fmt.Errorf("failed to parse request body: %w", err)
	}
	if parsed.Model == "" {
		return config.Provider{}, nil, fmt.Errorf(`request body is missing a "model" field`)
	}

	provider, ok := cfg.ModelRoutes[parsed.Model]
	if !ok {
		return config.Provider{}, nil, ErrUnknownModel
	}

	return provider, io.NopCloser(bytes.NewReader(raw)), nil
}