INSERT INTO hydra_oauth2_authentication_session
(id, authenticated_at, subject)
VALUES
('auth_session-0003', now(), 'subject-0003');

-- using the most lately added client as a foreign key
INSERT INTO hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id)
SELECT 'challenge-0003', 'verifier-0003', hydra_client.id, 'subject-0003', 'http://request/0003', true, 'requested_scope-0003_1', 'csrf-0003', now(), now(), '{"display": "display-0003"}', 'auth_session-0003'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge)
SELECT 'challenge-0003', 'verifier-0003', hydra_client.id, 'subject-0003', 'http://request/0003', true, 'requested_scope-0003_1', 'csrf-0003', now(), now(), '{"display": "display-0003"}', 'force_subject_id-0003', 'auth_session-0003', 'challenge-0003'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request_handled
(challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used)
VALUES
('challenge-0003', 'granted_scope-0003_1', true, 0003, '{}', now(), '{"session_access_token-0003": "0003"}', '{"session_id_token-0003": "0003"}', now(), true);

INSERT INTO hydra_oauth2_authentication_request_handled
(challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
('challenge-0003', 'subject-0003', true, 0003, '{}', 'acr-0003', now(), now(), true, 'force_subject_id-0003');

INSERT INTO hydra_oauth2_obfuscated_authentication_session (client_id, subject, subject_obfuscated)
SELECT hydra_client.id, 'subject-0003', 'subject_obfuscated-0003'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
