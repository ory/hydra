-- +migrate Up
ALTER TABLE hydra_oauth2_consent_request ADD username TEXT NULL ;
UPDATE hydra_oauth2_consent_request SET username='' ;
ALTER TABLE hydra_oauth2_consent_request MODIFY username TEXT NOT NULL ;

ALTER TABLE hydra_oauth2_authentication_request_handled ADD username TEXT NULL ;
UPDATE hydra_oauth2_authentication_request_handled SET username='' ;
ALTER TABLE hydra_oauth2_authentication_request_handled MODIFY username TEXT NOT NULL ;
-- +migrate Down
ALTER TABLE hydra_oauth2_consent_request DROP COLUMN username ;
ALTER TABLE hydra_oauth2_authentication_request_handled DROP COLUMN username ;