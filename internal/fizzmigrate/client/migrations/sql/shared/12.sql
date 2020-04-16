-- +migrate Up
ALTER TABLE hydra_client ADD created_at timestamp NOT NULL DEFAULT now();
ALTER TABLE hydra_client ADD updated_at timestamp NOT NULL DEFAULT now();

-- +migrate Down
ALTER TABLE hydra_client DROP COLUMN created_at;
ALTER TABLE hydra_client DROP COLUMN updated_at;
