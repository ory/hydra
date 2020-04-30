-- Fix performance issue of Admin API - Revoke Login Sessions
CREATE INDEX ON hydra_oauth2_authentication_session (subject);
