DROP TABLE IF EXISTS hydra_oauth2_device_auth_codes;

ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_challenge_id;
ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_code_request_id;
ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_verifier;
ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_csrf;
ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_was_used;
ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_handled_at;
ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_error;

ALTER TABLE hydra_client DROP COLUMN IF EXISTS device_authorization_grant_id_token_lifespan;
ALTER TABLE hydra_client DROP COLUMN IF EXISTS device_authorization_grant_access_token_lifespan;
ALTER TABLE hydra_client DROP COLUMN IF EXISTS device_authorization_grant_refresh_token_lifespan;
