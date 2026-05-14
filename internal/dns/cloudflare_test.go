package dns

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c, err := New("test-token", "test-zone-id")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestWithTimeout_addsDeadline(t *testing.T) {
	ctx, cancel := withTimeout(context.Background())
	defer cancel()

	if _, ok := ctx.Deadline(); !ok {
		t.Error("withTimeout: expected deadline to be set on a background context")
	}
}

func TestWithTimeout_preservesExisting(t *testing.T) {
	deadline := time.Now().Add(30 * time.Second)
	parent, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	ctx, cancel2 := withTimeout(parent)
	defer cancel2()

	got, ok := ctx.Deadline()
	if !ok {
		t.Fatal("deadline should still be set")
	}
	if !got.Equal(deadline) {
		t.Errorf("withTimeout changed deadline: got %v, want %v", got, deadline)
	}
}

// --- integration tests (require CF_TOKEN, CF_ZONE_ID, CF_DOMAIN env vars) ---

func skipIfNoCredentials(t *testing.T) (token, zoneID, domain string) {
	t.Helper()
	token = os.Getenv("CF_TOKEN")
	zoneID = os.Getenv("CF_ZONE_ID")
	domain = os.Getenv("CF_DOMAIN")
	if token == "" || zoneID == "" || domain == "" {
		t.Skip("CF_TOKEN, CF_ZONE_ID, CF_DOMAIN not set — skipping integration test")
	}
	return
}

func TestValidateToken_integration(t *testing.T) {
	token, zoneID, _ := skipIfNoCredentials(t)

	c, err := New(token, zoneID)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.ValidateToken(context.Background()); err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
}

func TestGetZoneID_integration(t *testing.T) {
	token, _, domain := skipIfNoCredentials(t)

	id, err := GetZoneID(context.Background(), token, domain)
	if err != nil {
		t.Fatalf("GetZoneID: %v", err)
	}
	if id == "" {
		t.Error("GetZoneID returned empty ID")
	}
}

func TestListRecords_integration(t *testing.T) {
	token, zoneID, domain := skipIfNoCredentials(t)

	c, err := New(token, zoneID)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	records, err := c.ListRecords(context.Background(), "*."+domain)
	if err != nil {
		t.Fatalf("ListRecords: %v", err)
	}
	t.Logf("found %d record(s) for *.%s", len(records), domain)
}
