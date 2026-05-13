// Package config handles loading and saving xpoz configuration.
package config

import (
	"fmt"
	"os"
)

// Config holds all xpoz configuration.
type Config struct {
	Domain     DomainConfig     `mapstructure:"domain"`
	Cloudflare CloudflareConfig `mapstructure:"cloudflare"`
	Proxy      ProxyConfig      `mapstructure:"proxy"`
	Subdomains SubdomainConfig  `mapstructure:"subdomains"`
}

// DomainConfig holds the user's domain name.
type DomainConfig struct {
	Name string `mapstructure:"name"`
}

// CloudflareConfig holds Cloudflare credentials.
// Token is never logged or included in error messages.
type CloudflareConfig struct {
	Token  string `mapstructure:"token"`
	ZoneID string `mapstructure:"zone_id"`
	Tunnel string `mapstructure:"tunnel"`
}

// ProxyConfig controls the internal reverse proxy.
type ProxyConfig struct {
	Port int `mapstructure:"port"`
}

// SubdomainConfig controls how subdomains are generated.
type SubdomainConfig struct {
	Style  string `mapstructure:"style"`  // hex | slug | word
	Length int    `mapstructure:"length"` // default 5
}

// Load reads config.toml from the xpoz data directory.
func Load() (*Config, error) {
	// TODO: implement in Phase 1
	return nil, fmt.Errorf("config: not implemented")
}

// Save writes cfg to disk with 0600 permissions.
func Save(cfg *Config) error {
	// TODO: implement in Phase 1
	return fmt.Errorf("config: not implemented")
}

// EnsureDir creates the xpoz data directory with 0700 permissions.
func EnsureDir() error {
	dir := Dir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("config: creating %s: %w", dir, err)
	}
	return nil
}
