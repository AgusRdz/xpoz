// Package dns manages Cloudflare DNS records for xpoz tunnels.
package dns

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

const defaultTimeout = 10 * time.Second

// Client wraps the Cloudflare API for xpoz DNS operations.
// The API token is held inside the cloudflare.API value and never surfaced in errors.
type Client struct {
	api    *cloudflare.API
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
func New(token, zoneID string) (*Client, error) {
	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return nil, fmt.Errorf("dns: creating client: %w", err)
	}
	return &Client{api: api, zoneID: zoneID}, nil
}

// ValidateToken verifies the token has DNS:Edit permission for the zone.
func (c *Client) ValidateToken(ctx context.Context) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	result, err := c.api.VerifyAPIToken(ctx)
	if err != nil {
		return fmt.Errorf("dns: verifying token: %w", err)
	}
	if result.Status != "active" {
		return fmt.Errorf("dns: token status is %q — check that the token has DNS:Edit permission for your zone", result.Status)
	}
	return nil
}

// GetZoneID resolves the Cloudflare zone ID for domain using token.
func GetZoneID(ctx context.Context, token, domain string) (string, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return "", fmt.Errorf("dns: creating client: %w", err)
	}
	zones, err := api.ListZones(ctx, domain)
	if err != nil {
		return "", fmt.Errorf("dns: listing zones for %q: %w", domain, err)
	}
	if len(zones) == 0 {
		return "", fmt.Errorf("dns: no Cloudflare zone found for %q — make sure the domain is added to your account", domain)
	}
	return zones[0].ID, nil
}

// EnsureWildcardCNAME creates *.domain → tunnelID.cfargotunnel.com (proxied)
// if the record does not already exist. Idempotent.
func (c *Client) EnsureWildcardCNAME(ctx context.Context, domain, tunnelID string) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	name := "*." + domain
	target := tunnelID + ".cfargotunnel.com"
	rc := cloudflare.ZoneIdentifier(c.zoneID)

	existing, _, err := c.api.ListDNSRecords(ctx, rc, cloudflare.ListDNSRecordsParams{
		Name: name,
		Type: "CNAME",
	})
	if err != nil {
		return fmt.Errorf("dns: listing records for %q: %w", name, err)
	}
	if len(existing) > 0 {
		return nil
	}

	proxied := true
	_, err = c.api.CreateDNSRecord(ctx, rc, cloudflare.CreateDNSRecordParams{
		Type:    "CNAME",
		Name:    name,
		Content: target,
		Proxied: &proxied,
		TTL:     1, // 1 = automatic TTL when proxied
	})
	if err != nil {
		return fmt.Errorf("dns: creating wildcard CNAME for %q: %w", domain, err)
	}
	return nil
}

// ListRecords returns all DNS records matching name in the zone.
func (c *Client) ListRecords(ctx context.Context, name string) ([]Record, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	rc := cloudflare.ZoneIdentifier(c.zoneID)
	cfRecords, _, err := c.api.ListDNSRecords(ctx, rc, cloudflare.ListDNSRecordsParams{Name: name})
	if err != nil {
		return nil, fmt.Errorf("dns: listing records for %q: %w", name, err)
	}

	records := make([]Record, len(cfRecords))
	for i, r := range cfRecords {
		records[i] = Record{
			ID:      r.ID,
			Type:    r.Type,
			Name:    r.Name,
			Content: r.Content,
			Proxied: r.Proxied != nil && *r.Proxied,
		}
	}
	return records, nil
}

// withTimeout wraps ctx with a 10s deadline if the caller did not set one.
func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, defaultTimeout)
}
