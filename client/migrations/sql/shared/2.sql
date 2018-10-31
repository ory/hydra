-- +migrate Up
ALTER TABLE hydra_client ADD client_secret_expires_at INTEGER NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE hydra_client DROP COLUMN client_secret_expires_at;
