DROP TABLE hydra_oauth2_device_auth_codes;

DROP INDEX hydra_oauth2_flow_device_challenge_idx;
ALTER TABLE hydra_oauth2_flow DROP COLUMN device_challenge_id;
ALTER TABLE hydra_oauth2_flow DROP COLUMN device_code_request_id;
ALTER TABLE hydra_oauth2_flow DROP COLUMN device_verifier;
ALTER TABLE hydra_oauth2_flow DROP COLUMN device_csrf;
ALTER TABLE hydra_oauth2_flow DROP COLUMN device_was_used;
ALTER TABLE hydra_oauth2_flow DROP COLUMN device_handled_at;
ALTER TABLE hydra_oauth2_flow DROP COLUMN device_error;

ALTER TABLE hydra_client DROP COLUMN device_authorization_grant_id_token_lifespan;
ALTER TABLE hydra_client DROP COLUMN device_authorization_grant_access_token_lifespan;
ALTER TABLE hydra_client DROP COLUMN device_authorization_grant_refresh_token_lifespan;
