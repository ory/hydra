ALTER TABLE hydra_oauth2_consent_request ADD amr TEXT NULL;
UPDATE hydra_oauth2_consent_request SET amr='';
ALTER TABLE hydra_oauth2_consent_request MODIFY amr TEXT NOT NULL;

ALTER TABLE hydra_oauth2_authentication_request_handled ADD amr TEXT NULL;
UPDATE hydra_oauth2_authentication_request_handled SET amr='';
ALTER TABLE hydra_oauth2_authentication_request_handled MODIFY amr TEXT NOT NULL;
