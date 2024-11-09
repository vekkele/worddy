-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  email citext UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sessions (
	token text PRIMARY KEY,
	data bytea NOT NULL,
	expiry timestamptz NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

-- +goose Down
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;