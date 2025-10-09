INSERT INTO hydra_oauth2_authentication_session
(id, authenticated_at, subject, remember)
VALUES
('auth_session-0014', '2022-02-15 22:20:21', 'subject-0014', false);

-- using the most lately added client as a foreign key
INSERT INTO hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, requested_at_audience, login_session_id)
SELECT 'challenge-0014', 'verifier-0014', hydra_client.id, 'subject-0014', 'http://request/0014', true, 'requested_scope-0014_1', 'csrf-0014', '2022-02-15 22:20:21', '2022-02-15 22:20:21', '{"display": "display-0014"}', 'requested_audience-0014_1', 'auth_session-0014'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr, context)
SELECT 'challenge-0014', 'verifier-0014', hydra_client.id, 'subject-0014', 'http://request/0014', true, 'requested_scope-0014_1', 'csrf-0014', '2022-02-15 22:20:21', '2022-02-15 22:20:21', '{"display": "display-0014"}', 'force_subject_id-0014', 'auth_session-0014', 'challenge-0014', 'requested_audience-0014_1', 'acr-0014', '{"context": "0014"}'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request_handled
(challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience, handled_at)
VALUES
('challenge-0014', 'granted_scope-0014_1', true, 0014, '{}', '2022-02-15 22:20:21', '{"session_access_token-0014": "0014"}', '{"session_id_token-0014": "0014"}', '2022-02-15 22:20:21', true, 'granted_audience-0014_1', '2022-02-15 22:20:21');

INSERT INTO hydra_oauth2_authentication_request_handled
(challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier, context)
VALUES
('challenge-0014', 'subject-0014', true, 0014, '{}', 'acr-0014', '2022-02-15 22:20:21', '2022-02-15 22:20:21', true, 'force_subject_id-0014', '{"context": "0014"}');

INSERT INTO hydra_oauth2_obfuscated_authentication_session (client_id, subject, subject_obfuscated)
SELECT hydra_client.id, 'subject-0014', 'subject_obfuscated-0014'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_logout_request (challenge, verifier, subject, sid, client_id, request_url, redir_url, was_used, accepted, rejected, rp_initiated)
SELECT 'challenge-0014', 'verifier-0014', 'subject-0014', 'session_id-0014', hydra_client.id, 'http://request/0014', 'http://post_logout/0014', true, true, false, true
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
