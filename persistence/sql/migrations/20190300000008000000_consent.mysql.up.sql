ALTER TABLE hydra_oauth2_authentication_request_handled ADD context TEXT NULL;
ALTER TABLE hydra_oauth2_consent_request ADD context TEXT NULL;

UPDATE hydra_oauth2_authentication_request_handled SET context='{}';
UPDATE hydra_oauth2_consent_request SET context='{}';

ALTER TABLE hydra_oauth2_authentication_request_handled MODIFY context TEXT NOT NULL;
ALTER TABLE hydra_oauth2_consent_request MODIFY context TEXT NOT NULL;

