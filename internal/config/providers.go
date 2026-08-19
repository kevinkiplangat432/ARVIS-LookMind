package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// Provider is a single upstream AI provider ARVIS can route to.
type Provider struct {
	Name    string   `yaml:"name"`
	BaseURL string   `yaml:"base_url"`
	APIKey  string   `yaml:"api_key"`
	Models  []string `yaml:"models"`
}

type providersFile struct {
	Providers []Provider `yaml:"providers"`
}

var envVarPattern = regexp.MustCompile(`\$\{([A-Z0-9_]+)\}`)

// LoadProviders reads and parses a providers YAML file, resolving any
// ${ENV_VAR} references in api_key fields against the current
// environment. Keys never live in the file itself — only the
// reference to where they come from does, which matters for a config
// file that's meant to be checked into version control.
func LoadProviders(path string) ([]Provider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read providers file: %w", err)
	}

	var pf providersFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("failed to parse providers file: %w", err)
	}

	if len(pf.Providers) == 0 {
		return nil, fmt.Errorf("no providers defined in %s", path)
	}
	if len(pf.Providers) > 20 {
		return nil, fmt.Errorf("providers file defines %d providers, ARVIS supports up to 20", len(pf.Providers))
	}

	seen := make(map[string]bool, len(pf.Providers))
	for i := range pf.Providers {
		p := &pf.Providers[i]
		if p.Name == "" {
			return nil, fmt.Errorf("provider at index %d is missing a name", i)
		}
		if seen[p.Name] {
			return nil, fmt.Errorf("duplicate provider name %q", p.Name)
		}
		seen[p.Name] = true

		p.APIKey = resolveEnvVars(p.APIKey)
		if p.APIKey == "" {
			return nil, fmt.Errorf("provider %q has no resolved api_key (check its env var is set)", p.Name)
		}
	}

	return pf.Providers, nil
}

func resolveEnvVars(s string) string {
	return envVarPattern.ReplaceAllStringFunc(s, func(match string) string {
		name := envVarPattern.FindStringSubmatch(match)[1]
		return os.Getenv(name)
	})
}