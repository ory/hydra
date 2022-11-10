ALTER TABLE hydra_oauth2_flow DROP CONSTRAINT hydra_oauth2_flow_login_session_id_fk;
ALTER TABLE hydra_oauth2_flow ADD CONSTRAINT hydra_oauth2_flow_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES hydra_oauth2_authentication_session(id) ON DELETE SET NULL;
ALTER TABLE hydra_oauth2_flow ALTER COLUMN login_session_id DROP DEFAULT;
