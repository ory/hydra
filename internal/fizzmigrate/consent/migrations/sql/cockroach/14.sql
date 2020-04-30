-- +migrate Up
ALTER TABLE hydra_oauth2_consent_request_handled ADD handled_at timestamp NULL;

-- +migrate Down
ALTER TABLE hydra_oauth2_consent_request_handled DROP COLUMN handled_at;
