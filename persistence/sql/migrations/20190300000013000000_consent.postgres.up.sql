-- Fix performance issue of Admin API - Revoke Login Sessions
CREATE INDEX hydra_oauth2_authentication_session_sub_idx ON hydra_oauth2_authentication_session (subject);
CREATE INDEX hydra_oauth2_authentication_request_login_session_id_idx ON hydra_oauth2_authentication_request (login_session_id);
CREATE INDEX hydra_oauth2_consent_request_login_session_id_idx ON hydra_oauth2_consent_request (login_session_id);
CREATE INDEX hydra_oauth2_consent_request_login_challenge_idx ON hydra_oauth2_consent_request (login_challenge);

-- Fix performance issue of Admin API - Revoke Consent Sessions
CREATE INDEX hydra_oauth2_logout_request_client_id_idx ON hydra_oauth2_logout_request (client_id);

