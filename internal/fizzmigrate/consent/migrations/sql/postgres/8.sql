-- +migrate Up
ALTER TABLE hydra_oauth2_authentication_request_handled ADD context TEXT NOT NULL DEFAULT '{}';
ALTER TABLE hydra_oauth2_consent_request ADD context TEXT NOT NULL DEFAULT '{}';

-- +migrate Down
ALTER TABLE hydra_oauth2_authentication_request_handled DROP COLUMN context;
ALTER TABLE hydra_oauth2_consent_request DROP COLUMN context;
