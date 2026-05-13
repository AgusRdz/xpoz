package store

// schema is the initial database schema applied on first open.
const schema = `
CREATE TABLE IF NOT EXISTS tunnels (
  name       TEXT PRIMARY KEY,
  subdomain  TEXT UNIQUE NOT NULL,
  full_url   TEXT NOT NULL,
  local_port INTEGER NOT NULL,
  active     BOOLEAN NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  last_used  DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tunnels_subdomain ON tunnels(subdomain);
`
