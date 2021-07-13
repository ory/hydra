-- Fix performance issue of Admin API - Revoke Login Sessions
DROP INDEX hydra_oauth2_authentication_session@hydra_oauth2_authentication_session_subject_idx;
