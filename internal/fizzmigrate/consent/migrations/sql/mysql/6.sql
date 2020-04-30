-- +migrate Up
ALTER TABLE hydra_oauth2_consent_request ADD acr TEXT NULL;
UPDATE hydra_oauth2_consent_request SET acr='';
ALTER TABLE hydra_oauth2_consent_request MODIFY acr TEXT NOT NULL;

-- +migrate Down
ALTER TABLE hydra_oauth2_consent_request DROP COLUMN acr;
