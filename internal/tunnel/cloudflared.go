// Package tunnel manages the cloudflared configuration file.
package tunnel

import (
	"fmt"
	"os"
	"runtime"
)

// Config represents the cloudflared config.yml structure.
type Config struct {
	Tunnel          string        `yaml:"tunnel"`
	CredentialsFile string        `yaml:"credentials-file"`
	Ingress         []IngressRule `yaml:"ingress"`
}

// IngressRule is one entry in the cloudflared ingress list.
type IngressRule struct {
	Hostname string `yaml:"hostname,omitempty"`
	Service  string `yaml:"service"`
}

// DefaultConfigPath returns the platform-specific cloudflared config path.
func DefaultConfigPath() string {
	if runtime.GOOS == "windows" {
		profile := os.Getenv("USERPROFILE")
		return profile + `\.cloudflared\config.yml`
	}
	home, _ := os.UserHomeDir()
	return home + "/.cloudflared/config.yml"
}

// Load reads and parses the cloudflared config file at path.
func Load(path string) (*Config, error) {
	// TODO: implement in Phase 1
	return nil, fmt.Errorf("tunnel: not implemented")
}

// Save writes cfg back to path using YAML marshaling.
// It preserves all existing ingress entries and enforces:
//   - xpoz wildcard entry is first
//   - catch-all service: http_status:404 is last
func Save(path string, cfg *Config) error {
	// TODO: implement in Phase 1
	return fmt.Errorf("tunnel: not implemented")
}

// EnsureWildcardIngress adds or updates the xpoz wildcard ingress rule in cfg.
// The rule routes *.domain to the internal proxy at proxyPort.
func EnsureWildcardIngress(cfg *Config, domain string, proxyPort int) {
	// TODO: implement in Phase 1
}
