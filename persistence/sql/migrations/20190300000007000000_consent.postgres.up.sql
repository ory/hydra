
-- This can be null when no previous login session exists, so let's remove default
ALTER TABLE hydra_oauth2_authentication_request ALTER COLUMN login_session_id DROP DEFAULT;

-- This can be null when no previous login session exists or if that session has been removed, so let's remove default
ALTER TABLE hydra_oauth2_consent_request ALTER COLUMN login_session_id DROP DEFAULT;

-- This can be null when the login_challenge was deleted (should not delete the consent itself)
ALTER TABLE hydra_oauth2_consent_request ALTER COLUMN login_challenge DROP DEFAULT;

-- Consent requests that point to an empty or invalid login request should set their login_challenge to NULL
UPDATE hydra_oauth2_consent_request SET login_challenge = NULL WHERE NOT EXISTS (
  SELECT 1 FROM hydra_oauth2_authentication_request WHERE hydra_oauth2_consent_request.login_challenge = hydra_oauth2_authentication_request.challenge
);

-- Consent requests that point to an empty or invalid login session should set their login_session_id to NULL
UPDATE hydra_oauth2_consent_request SET login_session_id = NULL WHERE NOT EXISTS (
  SELECT 1 FROM hydra_oauth2_authentication_session WHERE hydra_oauth2_consent_request.login_session_id = hydra_oauth2_authentication_session.id
);

-- Login requests that point to a login session that no longer exists (or was never set in the first place) should set that to NULL
UPDATE hydra_oauth2_authentication_request SET login_session_id = NULL WHERE NOT EXISTS (
  SELECT 1 FROM hydra_oauth2_authentication_session WHERE hydra_oauth2_authentication_request.login_session_id = hydra_oauth2_authentication_session.id
);

-- Login, consent, obfuscated sessions that point to a client which no longer exists must be deleted
DELETE FROM hydra_oauth2_authentication_request WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_authentication_request.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_consent_request WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_consent_request.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_obfuscated_authentication_session WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_obfuscated_authentication_session.client_id = hydra_client.id);

-- Handled login and consent requests which point to a consent/login request that no longer exists must be deleted
DELETE FROM hydra_oauth2_consent_request_handled WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_consent_request WHERE hydra_oauth2_consent_request_handled.challenge = hydra_oauth2_consent_request.challenge);
DELETE FROM hydra_oauth2_authentication_request_handled WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_consent_request WHERE hydra_oauth2_authentication_request_handled.challenge = hydra_oauth2_consent_request.challenge);

-- Actual indices

-- This caused #1209:
-- SET session_replication_role = replica;

-- Handled consent and authentication requests must cascade delete when their parent (the request itself) is removed
ALTER TABLE hydra_oauth2_consent_request_handled ADD CONSTRAINT hydra_oauth2_consent_request_handled_challenge_fk FOREIGN KEY (challenge) REFERENCES hydra_oauth2_consent_request(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_authentication_request_handled ADD CONSTRAINT hydra_oauth2_authentication_request_handled_challenge_fk FOREIGN KEY (challenge) REFERENCES hydra_oauth2_authentication_request(challenge) ON DELETE CASCADE;

-- Login, consent, obfuscated must be deleted when the oauth2 client is being deleted
ALTER TABLE hydra_oauth2_consent_request ADD CONSTRAINT hydra_oauth2_consent_request_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_authentication_request ADD CONSTRAINT hydra_oauth2_authentication_request_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_obfuscated_authentication_session ADD CONSTRAINT hydra_oauth2_obfuscated_authentication_session_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;

-- If a login session is removed, the associated login requests must be cascade deleted
ALTER TABLE hydra_oauth2_authentication_request ADD CONSTRAINT hydra_oauth2_authentication_request_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES hydra_oauth2_authentication_session(id) ON DELETE CASCADE;

-- But if a login session is removed the consent request should simply set it to NULL
ALTER TABLE hydra_oauth2_consent_request ADD CONSTRAINT hydra_oauth2_consent_request_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES hydra_oauth2_authentication_session(id) ON DELETE SET NULL;

-- It should also be set to null if the login request is deleted (because consent does not care about that)
ALTER TABLE hydra_oauth2_consent_request ADD CONSTRAINT hydra_oauth2_consent_request_login_challenge_fk FOREIGN KEY (login_challenge) REFERENCES hydra_oauth2_authentication_request(challenge) ON DELETE SET NULL;

-- This caused #1209:
-- SET session_replication_role = DEFAULT;
