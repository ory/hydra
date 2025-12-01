DROP TABLE hydra_oauth2_device_auth_codes;

ALTER TABLE hydra_oauth2_flow
  DROP COLUMN device_challenge_id,
  DROP COLUMN device_code_request_id,
  DROP COLUMN device_verifier,
  DROP COLUMN device_csrf,
  DROP COLUMN device_was_used,
  DROP COLUMN device_handled_at,
  DROP COLUMN device_error;

ALTER TABLE hydra_client
  DROP COLUMN device_authorization_grant_id_token_lifespan,
  DROP COLUMN device_authorization_grant_access_token_lifespan,
  DROP COLUMN device_authorization_grant_refresh_token_lifespan;
