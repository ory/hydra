-- +migrate Up
ALTER TABLE hydra_client ADD subject_type VARCHAR(15) NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE hydra_client DROP COLUMN subject_type;
