-- +migrate Up
ALTER TABLE hydra_oauth2_authentication_request MODIFY login_session_id VARCHAR(40) NULL;
ALTER TABLE hydra_oauth2_consent_request MODIFY login_session_id VARCHAR(40) NULL;
ALTER TABLE hydra_oauth2_consent_request ALTER login_session_id DROP DEFAULT;
ALTER TABLE hydra_oauth2_authentication_request ALTER login_session_id DROP DEFAULT;

UPDATE hydra_oauth2_authentication_request SET login_session_id = NULL WHERE login_session_id='';
UPDATE hydra_oauth2_consent_request SET login_session_id = NULL WHERE login_session_id='';

DELETE FROM hydra_oauth2_consent_request_handled WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_consent_request WHERE hydra_oauth2_consent_request_handled.challenge = hydra_oauth2_consent_request.challenge);
DELETE FROM hydra_oauth2_authentication_request_handled WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_consent_request WHERE hydra_oauth2_authentication_request_handled.challenge = hydra_oauth2_consent_request.challenge);

ALTER TABLE hydra_oauth2_consent_request_handled ADD CONSTRAINT hydra_oauth2_consent_request_handled_challenge_fk FOREIGN KEY (challenge) REFERENCES hydra_oauth2_consent_request(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_authentication_request_handled ADD CONSTRAINT hydra_oauth2_authentication_request_handled_challenge_fk FOREIGN KEY (challenge) REFERENCES hydra_oauth2_authentication_request(challenge) ON DELETE CASCADE;

DELETE FROM hydra_oauth2_consent_request WHERE login_challenge='';

DELETE FROM hydra_oauth2_authentication_request WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_authentication_request.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_authentication_request WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_authentication_session WHERE hydra_oauth2_authentication_request.login_session_id = hydra_oauth2_authentication_session.id);

DELETE FROM hydra_oauth2_consent_request WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_consent_request.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_consent_request WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_authentication_session WHERE hydra_oauth2_consent_request.login_session_id = hydra_oauth2_authentication_session.id);
DELETE FROM hydra_oauth2_consent_request WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_authentication_request WHERE hydra_oauth2_consent_request.login_challenge = hydra_oauth2_authentication_request.challenge);

DELETE FROM hydra_oauth2_obfuscated_authentication_session WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_obfuscated_authentication_session.client_id = hydra_client.id);

ALTER TABLE hydra_oauth2_consent_request ADD CONSTRAINT hydra_oauth2_consent_request_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_authentication_request ADD CONSTRAINT hydra_oauth2_authentication_request_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_obfuscated_authentication_session ADD CONSTRAINT hydra_oauth2_obfuscated_authentication_session_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;

ALTER TABLE hydra_oauth2_consent_request ADD CONSTRAINT hydra_oauth2_consent_request_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES hydra_oauth2_authentication_session(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_authentication_request ADD CONSTRAINT hydra_oauth2_authentication_request_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES hydra_oauth2_authentication_session(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_consent_request ADD CONSTRAINT hydra_oauth2_consent_request_login_challenge_fk FOREIGN KEY (login_challenge) REFERENCES hydra_oauth2_authentication_request(challenge) ON DELETE CASCADE;

-- +migrate Down
ALTER TABLE hydra_oauth2_consent_request_handled DROP FOREIGN KEY hydra_oauth2_consent_request_handled_challenge_fk;
ALTER TABLE hydra_oauth2_authentication_request_handled DROP FOREIGN KEY hydra_oauth2_authentication_request_handled_challenge_fk;

ALTER TABLE hydra_oauth2_consent_request DROP FOREIGN KEY hydra_oauth2_consent_request_client_id_fk;
ALTER TABLE hydra_oauth2_authentication_request DROP FOREIGN KEY hydra_oauth2_authentication_request_client_id_fk;
ALTER TABLE hydra_oauth2_obfuscated_authentication_session DROP FOREIGN KEY hydra_oauth2_obfuscated_authentication_session_client_id_fk;

ALTER TABLE hydra_oauth2_consent_request DROP FOREIGN KEY hydra_oauth2_consent_request_login_session_id_fk;
ALTER TABLE hydra_oauth2_authentication_request DROP FOREIGN KEY hydra_oauth2_authentication_request_login_session_id_fk;
ALTER TABLE hydra_oauth2_consent_request DROP FOREIGN KEY hydra_oauth2_consent_request_login_challenge_fk;
