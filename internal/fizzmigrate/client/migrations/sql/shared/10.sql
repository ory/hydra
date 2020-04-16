-- +migrate Up
ALTER TABLE hydra_client ADD audience TEXT;

-- +migrate Down
ALTER TABLE hydra_client DROP COLUMN audience;
