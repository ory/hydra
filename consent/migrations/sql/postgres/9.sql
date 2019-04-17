-- +migrate Up
ALTER TABLE hydra_oauth2_consent_request ADD consent_session_id TEXT NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE hydra_oauth2_consent_request DROP COLUMN consent_session_id;
