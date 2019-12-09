-- +migrate Up
-- Fix performance issue of Admin API - Revoke Login Sessions
CREATE INDEX ON hydra_oauth2_authentication_session (subject);

-- +migrate Down
-- Fix performance issue of Admin API - Revoke Login Sessions
DROP INDEX hydra_oauth2_authentication_session@hydra_oauth2_authentication_session_subject_idx;
