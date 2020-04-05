-- Fix performance issue of Admin API - Revoke Login Sessions
DROP INDEX hydra_oauth2_authentication_session_sub_idx ON hydra_oauth2_authentication_session;
-- The following 3 sets are for dropping indices for foreign keys. MySQL forbids me to drop indices on foreign keys: MySQL Cannot drop index needed in a foreign key constraint
ALTER TABLE hydra_oauth2_authentication_request DROP FOREIGN KEY hydra_oauth2_authentication_request_login_session_id_fk;
DROP INDEX hydra_oauth2_authentication_request_login_session_id_idx ON hydra_oauth2_authentication_request;
ALTER TABLE hydra_oauth2_authentication_request ADD CONSTRAINT hydra_oauth2_authentication_request_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES hydra_oauth2_authentication_session(id) ON DELETE CASCADE;

ALTER TABLE hydra_oauth2_consent_request DROP FOREIGN KEY hydra_oauth2_consent_request_login_session_id_fk;
DROP INDEX hydra_oauth2_consent_request_login_session_id_idx ON hydra_oauth2_consent_request;
ALTER TABLE hydra_oauth2_consent_request ADD CONSTRAINT hydra_oauth2_consent_request_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES hydra_oauth2_authentication_session(id) ON DELETE SET NULL;

ALTER TABLE hydra_oauth2_consent_request DROP FOREIGN KEY hydra_oauth2_consent_request_login_challenge_fk;
DROP INDEX hydra_oauth2_consent_request_login_challenge_idx ON hydra_oauth2_consent_request;
ALTER TABLE hydra_oauth2_consent_request ADD CONSTRAINT hydra_oauth2_consent_request_login_challenge_fk FOREIGN KEY (login_challenge) REFERENCES hydra_oauth2_authentication_request(challenge) ON DELETE SET NULL;

-- Fix performance issue of Admin API - Revoke Consent Sessions
ALTER TABLE hydra_oauth2_logout_request DROP FOREIGN KEY hydra_oauth2_logout_request_client_id_fk;
DROP INDEX hydra_oauth2_logout_request_client_id_idx ON hydra_oauth2_logout_request;
ALTER TABLE hydra_oauth2_logout_request ADD CONSTRAINT hydra_oauth2_logout_request_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
