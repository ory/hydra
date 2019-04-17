-- +migrate Up
ALTER TABLE hydra_oauth2_consent_request ADD consent_session_id TEXT NULL;
UPDATE hydra_oauth2_consent_request SET consent_session_id='';
ALTER TABLE hydra_oauth2_consent_request MODIFY consent_session_id TEXT NOT NULL;

-- +migrate Down
ALTER TABLE hydra_oauth2_consent_request DROP COLUMN consent_session_id;
