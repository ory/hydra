-- Fix performance issue of Admin API - Revoke Login Sessions
DROP INDEX hydra_oauth2_authentication_session_sub_idx;
DROP INDEX hydra_oauth2_authentication_request_login_session_id_idx;
DROP INDEX hydra_oauth2_consent_request_login_session_id_idx;
DROP INDEX hydra_oauth2_consent_request_login_challenge_idx;

-- Fix performance issue of Admin API - Revoke Consent Sessions
DROP INDEX hydra_oauth2_logout_request_client_id_idx;
