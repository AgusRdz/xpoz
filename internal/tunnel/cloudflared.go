// Package tunnel manages the cloudflared configuration file.
package tunnel

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Config represents the cloudflared config.yml structure.
// Unknown top-level keys (e.g. metrics, logfile) are captured in Extra
// so they survive a Load → EnsureWildcardIngress → Save round-trip.
type Config struct {
	Tunnel          string                 `yaml:"tunnel"`
	CredentialsFile string                 `yaml:"credentials-file"`
	Ingress         []IngressRule          `yaml:"ingress"`
	Extra           map[string]interface{} `yaml:",inline"`
}

// IngressRule is one entry in the cloudflared ingress list.
type IngressRule struct {
	Hostname string `yaml:"hostname,omitempty"`
	Service  string `yaml:"service"`
}

// DefaultConfigPath returns the platform-specific cloudflared config path.
func DefaultConfigPath() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE") + `\.cloudflared\config.yml`
	}
	home, _ := os.UserHomeDir()
	return home + "/.cloudflared/config.yml"
}

// Load reads and parses the cloudflared config file at path.
// Returns a wrapped os.ErrNotExist if the file is missing.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("tunnel: reading %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("tunnel: parsing %s: %w", path, err)
	}
	return &cfg, nil
}

// Save writes cfg to path using YAML marshaling with 0600 permissions.
// Fields captured in Config.Extra are preserved in the output.
func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("tunnel: marshaling config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("tunnel: writing %s: %w", path, err)
	}
	return nil
}

// EnsureWildcardIngress adds or updates the xpoz wildcard ingress rule in cfg.
// Guarantees: *.domain is first, catch-all http_status:404 is last,
// all other existing rules are preserved in their original order.
func EnsureWildcardIngress(cfg *Config, domain string, proxyPort int) {
	wildcardHost := "*." + domain
	service := fmt.Sprintf("http://localhost:%d", proxyPort)

	middle := make([]IngressRule, 0, len(cfg.Ingress))
	for _, rule := range cfg.Ingress {
		if rule.Hostname == wildcardHost || rule.Hostname == "" {
			continue
		}
		middle = append(middle, rule)
	}

	result := make([]IngressRule, 0, len(middle)+2)
	result = append(result, IngressRule{Hostname: wildcardHost, Service: service})
	result = append(result, middle...)
	result = append(result, IngressRule{Service: "http_status:404"})
	cfg.Ingress = result
}
