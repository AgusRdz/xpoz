// Package config handles loading and saving xpoz configuration.
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	DefaultProxyPort       = 2080
	DefaultSubdomainStyle  = "hex"
	DefaultSubdomainLength = 5
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

// Exists reports whether the config file is present on disk.
func Exists() bool {
	_, err := os.Stat(Path())
	return err == nil
}

// IsConfigured reports whether a valid, complete config file exists.
func IsConfigured() bool {
	cfg, err := Load()
	if err != nil {
		return false
	}
	return cfg.Validate() == nil
}

// Load reads config.toml from the xpoz data directory.
func Load() (*Config, error) {
	path := Path()
	if !Exists() {
		return nil, fmt.Errorf("config file not found at %s — run 'xpoz setup' to get started", path)
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config: reading %s: %w", path, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: parsing %s: %w", path, err)
	}

	cfg.applyDefaults()
	return &cfg, nil
}

// Save writes cfg to disk as TOML with 0600 permissions.
func Save(cfg *Config) error {
	if err := EnsureDir(); err != nil {
		return err
	}

	cfg.applyDefaults()

	content := fmt.Sprintf(tomlTemplate,
		cfg.Domain.Name,
		cfg.Cloudflare.Token,
		cfg.Cloudflare.ZoneID,
		cfg.Cloudflare.Tunnel,
		cfg.Proxy.Port,
		cfg.Subdomains.Style,
		cfg.Subdomains.Length,
	)

	path := Path()
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return fmt.Errorf("config: writing %s: %w", path, err)
	}

	return EnforcePermissions(path)
}

// Validate returns an error if any required field is missing.
func (c *Config) Validate() error {
	var missing []string
	if c.Domain.Name == "" {
		missing = append(missing, "domain.name")
	}
	if c.Cloudflare.Token == "" {
		missing = append(missing, "cloudflare.token")
	}
	if c.Cloudflare.ZoneID == "" {
		missing = append(missing, "cloudflare.zone_id")
	}
	if c.Cloudflare.Tunnel == "" {
		missing = append(missing, "cloudflare.tunnel")
	}
	if len(missing) > 0 {
		return fmt.Errorf("config: missing required fields: %s", strings.Join(missing, ", "))
	}
	return nil
}

// EnsureDir creates the xpoz data directory with 0700 permissions.
func EnsureDir() error {
	dir := Dir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("config: creating %s: %w", dir, err)
	}
	return nil
}

func (c *Config) applyDefaults() {
	if c.Proxy.Port == 0 {
		c.Proxy.Port = DefaultProxyPort
	}
	if c.Subdomains.Style == "" {
		c.Subdomains.Style = DefaultSubdomainStyle
	}
	if c.Subdomains.Length == 0 {
		c.Subdomains.Length = DefaultSubdomainLength
	}
}

// tomlTemplate is the canonical on-disk format.
// %q produces TOML-compatible double-quoted strings with proper escaping.
const tomlTemplate = `[domain]
name = %q

[cloudflare]
token   = %q
zone_id = %q
tunnel  = %q

[proxy]
port = %d

[subdomains]
style  = %q
length = %d
`
