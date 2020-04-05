-- +migrate Up
ALTER TABLE hydra_oauth2_consent_request ADD COLUMN username TEXT NOT NULL DEFAULT ('');
ALTER TABLE hydra_oauth2_authentication_request_handled ADD COLUMN username TEXT NOT NULL DEFAULT ('');
-- +migrate Down
ALTER TABLE hydra_oauth2_consent_request DROP COLUMN username;
ALTER TABLE hydra_oauth2_authentication_request_handled DROP COLUMN username;