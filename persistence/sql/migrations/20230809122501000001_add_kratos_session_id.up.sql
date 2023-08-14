ALTER TABLE hydra_oauth2_flow ADD COLUMN identity_provider_session_id VARCHAR(40);
ALTER TABLE hydra_oauth2_authentication_session ADD COLUMN identity_provider_session_id VARCHAR(40);