-- +migrate Up
ALTER TABLE hydra_client ADD metadata TEXT NOT NULL DEFAULT '{}';

-- +migrate Down
ALTER TABLE hydra_client DROP COLUMN metadata;
