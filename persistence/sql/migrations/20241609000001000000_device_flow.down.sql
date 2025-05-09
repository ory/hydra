ALTER TABLE hydra_oauth2_device_auth_codes DROP FOREIGN KEY IF EXISTS hydra_oauth2_device_auth_codes_challenge_id_fk;
ALTER TABLE hydra_oauth2_device_auth_codes DROP FOREIGN KEY IF EXISTS hydra_oauth2_device_auth_codes_client_id_fk;
ALTER TABLE hydra_oauth2_device_auth_codes DROP FOREIGN KEY IF EXISTS hydra_oauth2_device_auth_codes_nid_fk_idx;

DROP TABLE IF EXISTS hydra_oauth2_device_auth_codes;

ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_challenge_id;
ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_code_request_id;
ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_verifier;
ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_csrf;
ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_was_used;
ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_handled_at;
ALTER TABLE hydra_oauth2_flow DROP COLUMN IF EXISTS device_error;
