ALTER TABLE hydra_oauth2_consent_request ADD login_session_id VARCHAR(40) NULL DEFAULT '';
ALTER TABLE hydra_oauth2_consent_request ADD login_challenge VARCHAR(40) NULL DEFAULT '';
ALTER TABLE hydra_oauth2_authentication_request ADD login_session_id VARCHAR(40) NULL DEFAULT '';
