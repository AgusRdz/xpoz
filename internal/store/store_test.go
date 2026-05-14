package store

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func newTestStore(t *testing.T) Store {
	t.Helper()
	s, err := New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func makeTunnel(name, subdomain string, port int) *Tunnel {
	now := time.Now().UTC().Truncate(time.Second)
	return &Tunnel{
		Name:      name,
		Subdomain: subdomain,
		FullURL:   "https://" + subdomain + ".example.com",
		LocalPort: port,
		Active:    false,
		CreatedAt: now,
		LastUsed:  now,
	}
}

func TestNew(t *testing.T) {
	s := newTestStore(t)
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestUpsertAndGet(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	want := makeTunnel("my-app", "a3f9c", 3000)
	want.Active = true

	if err := s.Upsert(ctx, want); err != nil {
		t.Fatalf("Upsert: %v", err)
	}

	got, err := s.Get(ctx, "my-app")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if got.Name != want.Name {
		t.Errorf("Name: got %q, want %q", got.Name, want.Name)
	}
	if got.Subdomain != want.Subdomain {
		t.Errorf("Subdomain: got %q, want %q", got.Subdomain, want.Subdomain)
	}
	if got.FullURL != want.FullURL {
		t.Errorf("FullURL: got %q, want %q", got.FullURL, want.FullURL)
	}
	if got.LocalPort != want.LocalPort {
		t.Errorf("LocalPort: got %d, want %d", got.LocalPort, want.LocalPort)
	}
	if got.Active != want.Active {
		t.Errorf("Active: got %v, want %v", got.Active, want.Active)
	}
	if !got.CreatedAt.Equal(want.CreatedAt) {
		t.Errorf("CreatedAt: got %v, want %v", got.CreatedAt, want.CreatedAt)
	}
}

func TestUpsert_updatesExisting(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	original := makeTunnel("my-app", "a3f9c", 3000)
	if err := s.Upsert(ctx, original); err != nil {
		t.Fatalf("Upsert original: %v", err)
	}

	updated := makeTunnel("my-app", "b7e2d", 4000)
	updated.CreatedAt = original.CreatedAt.Add(time.Hour) // should NOT be persisted
	if err := s.Upsert(ctx, updated); err != nil {
		t.Fatalf("Upsert updated: %v", err)
	}

	got, err := s.Get(ctx, "my-app")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if got.Subdomain != "b7e2d" {
		t.Errorf("Subdomain not updated: got %q", got.Subdomain)
	}
	if got.LocalPort != 4000 {
		t.Errorf("LocalPort not updated: got %d", got.LocalPort)
	}
	// created_at must be preserved from the original insert
	if !got.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("CreatedAt changed on update: got %v, want %v", got.CreatedAt, original.CreatedAt)
	}
}

func TestGet_notFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.Get(ctx, "nonexistent")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetBySubdomain(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	want := makeTunnel("my-app", "a3f9c", 3000)
	if err := s.Upsert(ctx, want); err != nil {
		t.Fatalf("Upsert: %v", err)
	}

	got, err := s.GetBySubdomain(ctx, "a3f9c")
	if err != nil {
		t.Fatalf("GetBySubdomain: %v", err)
	}
	if got.Name != "my-app" {
		t.Errorf("Name: got %q, want %q", got.Name, "my-app")
	}
}

func TestGetBySubdomain_notFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.GetBySubdomain(ctx, "zzzzz")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	tunnels := []*Tunnel{
		makeTunnel("app-a", "aaaaa", 3000),
		makeTunnel("app-b", "bbbbb", 4000),
		makeTunnel("app-c", "ccccc", 5000),
	}
	// insert with different created_at so ordering is deterministic
	for i, tun := range tunnels {
		tun.CreatedAt = time.Now().UTC().Add(time.Duration(i) * time.Second).Truncate(time.Second)
		tun.LastUsed = tun.CreatedAt
		if err := s.Upsert(ctx, tun); err != nil {
			t.Fatalf("Upsert %q: %v", tun.Name, err)
		}
	}

	got, err := s.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("List count: got %d, want 3", len(got))
	}
	// newest first
	if got[0].Name != "app-c" {
		t.Errorf("first item: got %q, want %q", got[0].Name, "app-c")
	}
	if got[2].Name != "app-a" {
		t.Errorf("last item: got %q, want %q", got[2].Name, "app-a")
	}
}

func TestList_empty(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	got, err := s.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if got == nil {
		t.Error("List returned nil, want empty slice")
	}
	if len(got) != 0 {
		t.Errorf("List count: got %d, want 0", len(got))
	}
}

func TestDelete(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	if err := s.Upsert(ctx, makeTunnel("my-app", "a3f9c", 3000)); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	if err := s.Delete(ctx, "my-app"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Get(ctx, "my-app")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDelete_notFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.Delete(ctx, "nonexistent")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSetActive(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	if err := s.Upsert(ctx, makeTunnel("my-app", "a3f9c", 3000)); err != nil {
		t.Fatalf("Upsert: %v", err)
	}

	if err := s.SetActive(ctx, "my-app", true); err != nil {
		t.Fatalf("SetActive true: %v", err)
	}
	got, _ := s.Get(ctx, "my-app")
	if !got.Active {
		t.Error("Active should be true")
	}

	if err := s.SetActive(ctx, "my-app", false); err != nil {
		t.Fatalf("SetActive false: %v", err)
	}
	got, _ = s.Get(ctx, "my-app")
	if got.Active {
		t.Error("Active should be false")
	}
}

func TestSetActive_notFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.SetActive(ctx, "nonexistent", true)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
