package tunnel

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writing test config: %v", err)
	}
	return path
}

func TestLoad(t *testing.T) {
	path := writeConfig(t, `
tunnel: abc123
credentials-file: /home/user/.cloudflared/abc123.json
ingress:
  - hostname: app.example.com
    service: http://localhost:3000
  - service: http_status:404
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Tunnel != "abc123" {
		t.Errorf("Tunnel: got %q, want %q", cfg.Tunnel, "abc123")
	}
	if len(cfg.Ingress) != 2 {
		t.Fatalf("Ingress len: got %d, want 2", len(cfg.Ingress))
	}
	if cfg.Ingress[0].Hostname != "app.example.com" {
		t.Errorf("Ingress[0].Hostname: got %q", cfg.Ingress[0].Hostname)
	}
}

func TestLoad_missing(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yml"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestLoad_invalidYAML(t *testing.T) {
	path := writeConfig(t, "{{{")
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestLoad_preservesUnknownFields(t *testing.T) {
	path := writeConfig(t, `
tunnel: abc123
metrics: localhost:3001
logfile: /var/log/cloudflared.log
ingress:
  - service: http_status:404
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Extra["metrics"] != "localhost:3001" {
		t.Errorf("Extra[metrics]: got %v", cfg.Extra["metrics"])
	}
	if cfg.Extra["logfile"] != "/var/log/cloudflared.log" {
		t.Errorf("Extra[logfile]: got %v", cfg.Extra["logfile"])
	}
}

func TestSaveLoad_roundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	original := &Config{
		Tunnel:          "abc123",
		CredentialsFile: "/home/user/.cloudflared/abc123.json",
		Ingress: []IngressRule{
			{Hostname: "*.example.com", Service: "http://localhost:2080"},
			{Service: "http_status:404"},
		},
	}
	if err := Save(path, original); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load after Save: %v", err)
	}
	if got.Tunnel != original.Tunnel {
		t.Errorf("Tunnel: got %q, want %q", got.Tunnel, original.Tunnel)
	}
	if len(got.Ingress) != len(original.Ingress) {
		t.Fatalf("Ingress len: got %d, want %d", len(got.Ingress), len(original.Ingress))
	}
}

func TestEnsureWildcardIngress_addsNew(t *testing.T) {
	cfg := &Config{
		Ingress: []IngressRule{
			{Hostname: "app.example.com", Service: "http://localhost:3000"},
			{Service: "http_status:404"},
		},
	}
	EnsureWildcardIngress(cfg, "example.com", 2080)

	if len(cfg.Ingress) != 3 {
		t.Fatalf("Ingress len: got %d, want 3", len(cfg.Ingress))
	}
	if cfg.Ingress[0].Hostname != "*.example.com" {
		t.Errorf("first rule hostname: got %q, want *.example.com", cfg.Ingress[0].Hostname)
	}
	if cfg.Ingress[0].Service != "http://localhost:2080" {
		t.Errorf("wildcard service: got %q", cfg.Ingress[0].Service)
	}
}

func TestEnsureWildcardIngress_updatesExisting(t *testing.T) {
	cfg := &Config{
		Ingress: []IngressRule{
			{Hostname: "*.example.com", Service: "http://localhost:9999"},
			{Hostname: "app.example.com", Service: "http://localhost:3000"},
			{Service: "http_status:404"},
		},
	}
	EnsureWildcardIngress(cfg, "example.com", 2080)

	if len(cfg.Ingress) != 3 {
		t.Fatalf("Ingress len: got %d, want 3 (no duplicate wildcard)", len(cfg.Ingress))
	}
	if cfg.Ingress[0].Service != "http://localhost:2080" {
		t.Errorf("wildcard service not updated: got %q", cfg.Ingress[0].Service)
	}
}

func TestEnsureWildcardIngress_catchAllIsLast(t *testing.T) {
	cfg := &Config{}
	EnsureWildcardIngress(cfg, "example.com", 2080)

	last := cfg.Ingress[len(cfg.Ingress)-1]
	if last.Hostname != "" || last.Service != "http_status:404" {
		t.Errorf("last rule should be catch-all, got %+v", last)
	}
}

func TestEnsureWildcardIngress_preservesOtherRules(t *testing.T) {
	cfg := &Config{
		Ingress: []IngressRule{
			{Hostname: "api.example.com", Service: "http://localhost:4000"},
			{Hostname: "app.example.com", Service: "http://localhost:3000"},
			{Service: "http_status:404"},
		},
	}
	EnsureWildcardIngress(cfg, "example.com", 2080)

	var found []string
	for _, r := range cfg.Ingress {
		if r.Hostname != "" && r.Hostname != "*.example.com" {
			found = append(found, r.Hostname)
		}
	}
	if len(found) != 2 {
		t.Errorf("expected 2 preserved rules, got %v", found)
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	if path == "" {
		t.Error("DefaultConfigPath returned empty string")
	}
	if !strings.Contains(path, "cloudflared") {
		t.Errorf("DefaultConfigPath %q does not contain 'cloudflared'", path)
	}
}
