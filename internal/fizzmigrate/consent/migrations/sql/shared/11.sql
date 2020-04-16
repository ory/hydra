-- +migrate Up
ALTER TABLE hydra_oauth2_authentication_session ADD remember bool NOT NULL DEFAULT FALSE;

UPDATE hydra_oauth2_authentication_session SET remember=TRUE;

-- +migrate Down
ALTER TABLE hydra_oauth2_authentication_session DROP COLUMN remember;
