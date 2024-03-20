-- +goose Up
CREATE TABLE sessions (
	token  TEXT PRIMARY KEY,
	data   BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_sessions_expiry ON sessions (expiry);



