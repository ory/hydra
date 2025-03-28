CREATE INDEX IF NOT EXISTS hydra_oauth2_flow_nid_sid_subject_idx ON hydra_oauth2_flow (nid, login_session_id, subject) STORING (client_id) WHERE login_session_id IS NOT NULL;
