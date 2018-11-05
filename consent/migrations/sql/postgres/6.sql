-- +migrate Up
ALTER TABLE hydra_oauth2_consent_request ADD acr TEXT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE hydra_oauth2_consent_request DROP COLUMN acr;
