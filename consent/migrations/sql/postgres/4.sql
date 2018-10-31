-- +migrate Up
ALTER TABLE hydra_oauth2_consent_request ADD requested_at_audience TEXT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_authentication_request ADD requested_at_audience TEXT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_consent_request_handled ADD granted_at_audience TEXT NULL DEFAULT '';
				
-- +migrate Down
ALTER TABLE hydra_oauth2_consent_request DROP COLUMN requested_at_audience;
ALTER TABLE hydra_oauth2_authentication_request DROP COLUMN requested_at_audience;
ALTER TABLE hydra_oauth2_consent_request_handled DROP COLUMN granted_at_audience;
