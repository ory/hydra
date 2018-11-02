-- +migrate Up
ALTER TABLE hydra_oauth2_consent_request ADD login_session_id VARCHAR(40) NULL DEFAULT '';
ALTER TABLE hydra_oauth2_consent_request ADD login_challenge VARCHAR(40) NULL DEFAULT '';
ALTER TABLE hydra_oauth2_authentication_request ADD login_session_id VARCHAR(40) NULL DEFAULT '';
				
-- +migrate Down
ALTER TABLE hydra_oauth2_consent_request DROP COLUMN login_session_id;
ALTER TABLE hydra_oauth2_consent_request DROP COLUMN login_challenge;
ALTER TABLE hydra_oauth2_authentication_request DROP COLUMN login_session_id;
