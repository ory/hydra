INSERT INTO hydra_oauth2_authentication_session
(id, authenticated_at, subject)
VALUES
('auth_session-0005', now(), 'subject-0005');

-- using the most lately added client as a foreign key
INSERT INTO hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, requested_at_audience, login_session_id)
SELECT 'challenge-0005', 'verifier-0005', hydra_client.id, 'subject-0005', 'http://request/0005', true, 'requested_scope-0005_1', 'csrf-0005', now(), now(), '{"display": "display-0005"}', 'requested_audience-0005_1', 'auth_session-0005'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience)
SELECT 'challenge-0005', 'verifier-0005', hydra_client.id, 'subject-0005', 'http://request/0005', true, 'requested_scope-0005_1', 'csrf-0005', now(), now(), '{"display": "display-0005"}', 'force_subject_id-0005', 'auth_session-0005', 'challenge-0005', 'requested_audience-0005_1'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request_handled
(challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
('challenge-0005', 'granted_scope-0005_1', true, 0005, '{}', now(), '{"session_access_token-0005": "0005"}', '{"session_id_token-0005": "0005"}', now(), true, 'granted_audience-0005_1');

INSERT INTO hydra_oauth2_authentication_request_handled
(challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
('challenge-0005', 'subject-0005', true, 0005, '{}', 'acr-0005', now(), now(), true, 'force_subject_id-0005');

INSERT INTO hydra_oauth2_obfuscated_authentication_session (client_id, subject, subject_obfuscated)
SELECT hydra_client.id, 'subject-0005', 'subject_obfuscated-0005'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
