package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite" // registers the "sqlite" driver
)

const timeLayout = time.RFC3339

type sqliteStore struct {
	db *sql.DB
}

// rowScanner is satisfied by both *sql.Row and *sql.Rows.
type rowScanner interface {
	Scan(dest ...any) error
}

// New opens (or creates) the SQLite database at path and runs schema migrations.
func New(path string) (Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("store: opening %s: %w", path, err)
	}

	pragmas := []string{
		`PRAGMA journal_mode=WAL`,
		`PRAGMA busy_timeout=5000`,
		`PRAGMA foreign_keys=ON`,
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			db.Close()
			return nil, fmt.Errorf("store: %s: %w", p, err)
		}
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("store: applying schema: %w", err)
	}

	return &sqliteStore{db: db}, nil
}

func (s *sqliteStore) Get(ctx context.Context, name string) (*Tunnel, error) {
	const q = `SELECT name, subdomain, full_url, local_port, active, created_at, last_used
	           FROM tunnels WHERE name = ?`
	t, err := scanTunnel(s.db.QueryRowContext(ctx, q, name))
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("store: get %q: %w", name, err)
	}
	return t, nil
}

func (s *sqliteStore) GetBySubdomain(ctx context.Context, subdomain string) (*Tunnel, error) {
	const q = `SELECT name, subdomain, full_url, local_port, active, created_at, last_used
	           FROM tunnels WHERE subdomain = ?`
	t, err := scanTunnel(s.db.QueryRowContext(ctx, q, subdomain))
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("store: get by subdomain %q: %w", subdomain, err)
	}
	return t, nil
}

func (s *sqliteStore) List(ctx context.Context) ([]*Tunnel, error) {
	const q = `SELECT name, subdomain, full_url, local_port, active, created_at, last_used
	           FROM tunnels ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("store: list: %w", err)
	}
	defer rows.Close()

	tunnels := make([]*Tunnel, 0)
	for rows.Next() {
		t, err := scanTunnel(rows)
		if err != nil {
			return nil, fmt.Errorf("store: list scan: %w", err)
		}
		tunnels = append(tunnels, t)
	}
	return tunnels, rows.Err()
}

func (s *sqliteStore) Upsert(ctx context.Context, t *Tunnel) error {
	const q = `
INSERT INTO tunnels (name, subdomain, full_url, local_port, active, created_at, last_used)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(name) DO UPDATE SET
    subdomain  = excluded.subdomain,
    full_url   = excluded.full_url,
    local_port = excluded.local_port,
    active     = excluded.active,
    last_used  = excluded.last_used`

	_, err := s.db.ExecContext(ctx, q,
		t.Name, t.Subdomain, t.FullURL, t.LocalPort,
		boolToInt(t.Active),
		t.CreatedAt.UTC().Format(timeLayout),
		t.LastUsed.UTC().Format(timeLayout),
	)
	if err != nil {
		return fmt.Errorf("store: upsert %q: %w", t.Name, err)
	}
	return nil
}

func (s *sqliteStore) Delete(ctx context.Context, name string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM tunnels WHERE name = ?`, name)
	if err != nil {
		return fmt.Errorf("store: delete %q: %w", name, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *sqliteStore) SetActive(ctx context.Context, name string, active bool) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE tunnels SET active = ? WHERE name = ?`, boolToInt(active), name)
	if err != nil {
		return fmt.Errorf("store: set active %q: %w", name, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *sqliteStore) Close() error {
	return s.db.Close()
}

func scanTunnel(row rowScanner) (*Tunnel, error) {
	var t Tunnel
	var active int
	var createdAt, lastUsed string

	if err := row.Scan(&t.Name, &t.Subdomain, &t.FullURL, &t.LocalPort,
		&active, &createdAt, &lastUsed); err != nil {
		return nil, err
	}

	t.Active = active != 0
	t.CreatedAt, _ = time.Parse(timeLayout, createdAt)
	t.LastUsed, _ = time.Parse(timeLayout, lastUsed)
	return &t, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
