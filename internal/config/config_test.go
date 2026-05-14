package config

import (
	"os"
	"runtime"
	"testing"
)

// setupTempDir redirects Dir() to a temporary directory for the duration of the test.
func setupTempDir(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	if runtime.GOOS == "windows" {
		t.Setenv("APPDATA", dir)
	} else {
		t.Setenv("HOME", dir)
	}
}

func validConfig() *Config {
	return &Config{
		Domain:     DomainConfig{Name: "example.com"},
		Cloudflare: CloudflareConfig{Token: "test-token", ZoneID: "test-zone-id", Tunnel: "test-tunnel-id"},
		Proxy:      ProxyConfig{Port: 2080},
		Subdomains: SubdomainConfig{Style: "hex", Length: 5},
	}
}

func TestSaveLoad(t *testing.T) {
	setupTempDir(t)

	want := validConfig()
	if err := Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got.Domain.Name != want.Domain.Name {
		t.Errorf("Domain.Name: got %q, want %q", got.Domain.Name, want.Domain.Name)
	}
	if got.Cloudflare.Token != want.Cloudflare.Token {
		t.Errorf("Cloudflare.Token: got %q, want %q", got.Cloudflare.Token, want.Cloudflare.Token)
	}
	if got.Cloudflare.ZoneID != want.Cloudflare.ZoneID {
		t.Errorf("Cloudflare.ZoneID: got %q, want %q", got.Cloudflare.ZoneID, want.Cloudflare.ZoneID)
	}
	if got.Cloudflare.Tunnel != want.Cloudflare.Tunnel {
		t.Errorf("Cloudflare.Tunnel: got %q, want %q", got.Cloudflare.Tunnel, want.Cloudflare.Tunnel)
	}
	if got.Proxy.Port != want.Proxy.Port {
		t.Errorf("Proxy.Port: got %d, want %d", got.Proxy.Port, want.Proxy.Port)
	}
	if got.Subdomains.Style != want.Subdomains.Style {
		t.Errorf("Subdomains.Style: got %q, want %q", got.Subdomains.Style, want.Subdomains.Style)
	}
	if got.Subdomains.Length != want.Subdomains.Length {
		t.Errorf("Subdomains.Length: got %d, want %d", got.Subdomains.Length, want.Subdomains.Length)
	}
}

func TestSaveLoad_specialChars(t *testing.T) {
	setupTempDir(t)

	want := &Config{
		Domain:     DomainConfig{Name: "my-domain.com"},
		Cloudflare: CloudflareConfig{Token: `tok"en\with'special`, ZoneID: "zone", Tunnel: "tunnel"},
	}
	if err := Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Cloudflare.Token != want.Cloudflare.Token {
		t.Errorf("Token round-trip: got %q, want %q", got.Cloudflare.Token, want.Cloudflare.Token)
	}
}

func TestLoad_missing(t *testing.T) {
	setupTempDir(t)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestLoad_invalidTOML(t *testing.T) {
	setupTempDir(t)
	if err := EnsureDir(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(Path(), []byte("not valid toml :::"), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid TOML, got nil")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name:    "valid",
			cfg:     validConfig(),
			wantErr: false,
		},
		{
			name:    "missing domain",
			cfg:     &Config{Cloudflare: CloudflareConfig{Token: "t", ZoneID: "z", Tunnel: "tu"}},
			wantErr: true,
		},
		{
			name:    "missing token",
			cfg:     &Config{Domain: DomainConfig{Name: "x.com"}, Cloudflare: CloudflareConfig{ZoneID: "z", Tunnel: "tu"}},
			wantErr: true,
		},
		{
			name:    "missing zone_id",
			cfg:     &Config{Domain: DomainConfig{Name: "x.com"}, Cloudflare: CloudflareConfig{Token: "t", Tunnel: "tu"}},
			wantErr: true,
		},
		{
			name:    "missing tunnel",
			cfg:     &Config{Domain: DomainConfig{Name: "x.com"}, Cloudflare: CloudflareConfig{Token: "t", ZoneID: "z"}},
			wantErr: true,
		},
		{
			name:    "all empty",
			cfg:     &Config{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestApplyDefaults(t *testing.T) {
	cfg := &Config{}
	cfg.applyDefaults()

	if cfg.Proxy.Port != DefaultProxyPort {
		t.Errorf("Proxy.Port: got %d, want %d", cfg.Proxy.Port, DefaultProxyPort)
	}
	if cfg.Subdomains.Style != DefaultSubdomainStyle {
		t.Errorf("Subdomains.Style: got %q, want %q", cfg.Subdomains.Style, DefaultSubdomainStyle)
	}
	if cfg.Subdomains.Length != DefaultSubdomainLength {
		t.Errorf("Subdomains.Length: got %d, want %d", cfg.Subdomains.Length, DefaultSubdomainLength)
	}
}

func TestApplyDefaults_preservesExisting(t *testing.T) {
	cfg := &Config{
		Proxy:      ProxyConfig{Port: 9090},
		Subdomains: SubdomainConfig{Style: "slug", Length: 8},
	}
	cfg.applyDefaults()

	if cfg.Proxy.Port != 9090 {
		t.Errorf("Proxy.Port should not be overwritten: got %d", cfg.Proxy.Port)
	}
	if cfg.Subdomains.Style != "slug" {
		t.Errorf("Subdomains.Style should not be overwritten: got %q", cfg.Subdomains.Style)
	}
}

func TestExists(t *testing.T) {
	setupTempDir(t)

	if Exists() {
		t.Error("Exists() = true before any Save, want false")
	}

	if err := Save(validConfig()); err != nil {
		t.Fatal(err)
	}

	if !Exists() {
		t.Error("Exists() = false after Save, want true")
	}
}

func TestIsConfigured(t *testing.T) {
	setupTempDir(t)

	if IsConfigured() {
		t.Error("IsConfigured() = true with no config file, want false")
	}

	if err := Save(validConfig()); err != nil {
		t.Fatal(err)
	}

	if !IsConfigured() {
		t.Error("IsConfigured() = false after saving valid config, want true")
	}
}

func TestSave_permissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission bits not enforced on Windows")
	}
	setupTempDir(t)

	if err := Save(validConfig()); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(Path())
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file permissions: got %04o, want 0600", perm)
	}
}

func TestEnsureDir_permissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission bits not enforced on Windows")
	}
	setupTempDir(t)

	if err := EnsureDir(); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(Dir())
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0700 {
		t.Errorf("dir permissions: got %04o, want 0700", perm)
	}
}
