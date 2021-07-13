-- +migrate Up
ALTER TABLE hydra_jwk ADD created_at TIMESTAMP NOT NULL DEFAULT NOW();

-- +migrate Down
ALTER TABLE hydra_jwk DROP COLUMN created_at;
