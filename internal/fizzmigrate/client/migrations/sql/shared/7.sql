-- +migrate Up
ALTER TABLE hydra_client ADD allowed_cors_origins TEXT;

-- +migrate Down
ALTER TABLE hydra_client DROP COLUMN allowed_cors_origins;
