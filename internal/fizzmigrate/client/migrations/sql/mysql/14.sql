-- +migrate Up
ALTER TABLE hydra_client ADD metadata TEXT NULL;

UPDATE hydra_client SET metadata='{}';

ALTER TABLE hydra_client MODIFY metadata TEXT NOT NULL;

-- +migrate Down
ALTER TABLE hydra_client DROP COLUMN metadata;
