CREATE INDEX CONCURRENTLY IF NOT EXISTS hydra_oauth2_flow_login_session_subject_idx ON hydra_oauth2_flow (nid, login_session_id, subject) INCLUDE (client_id) WHERE login_session_id IS NOT NULL;
