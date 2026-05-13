// Package dns manages Cloudflare DNS records for xpoz tunnels.
package dns

import (
	"context"
	"fmt"
)

// Client wraps the Cloudflare API for xpoz DNS operations.
// Token is stored but never logged or included in error messages.
type Client struct {
	token  string
	zoneID string
}

// Record is a simplified Cloudflare DNS record.
type Record struct {
	ID      string
	Type    string
	Name    string
	Content string
	Proxied bool
}

// New creates a Client for the given token and zone.
func New(token, zoneID string) *Client {
	return &Client{token: token, zoneID: zoneID}
}

// ValidateToken verifies the token has DNS:Edit permission for the zone.
func (c *Client) ValidateToken(ctx context.Context) error {
	// TODO: implement in Phase 1
	return fmt.Errorf("dns: not implemented")
}

// GetZoneID resolves the zone ID for domain using token.
func GetZoneID(ctx context.Context, token, domain string) (string, error) {
	// TODO: implement in Phase 1
	return "", fmt.Errorf("dns: not implemented")
}

// EnsureWildcardCNAME creates *.domain → tunnelID.cfargotunnel.com (proxied)
// if the record does not already exist.
func (c *Client) EnsureWildcardCNAME(ctx context.Context, domain, tunnelID string) error {
	// TODO: implement in Phase 1
	return fmt.Errorf("dns: not implemented")
}

// ListRecords returns existing DNS records matching name.
func (c *Client) ListRecords(ctx context.Context, name string) ([]Record, error) {
	// TODO: implement in Phase 1
	return nil, fmt.Errorf("dns: not implemented")
}
