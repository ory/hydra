INSERT INTO hydra_oauth2_authentication_session
(id, authenticated_at, subject)
VALUES
('auth_session-0010', '2022-02-15 22:20:21', 'subject-0010');

-- using the most lately added client as a foreign key
INSERT INTO hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, requested_at_audience, login_session_id)
SELECT 'challenge-0010', 'verifier-0010', hydra_client.id, 'subject-0010', 'http://request/0010', true, 'requested_scope-0010_1', 'csrf-0010', '2022-02-15 22:20:21', '2022-02-15 22:20:21', '{"display": "display-0010"}', 'requested_audience-0010_1', 'auth_session-0010'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr, context)
SELECT 'challenge-0010', 'verifier-0010', hydra_client.id, 'subject-0010', 'http://request/0010', true, 'requested_scope-0010_1', 'csrf-0010', '2022-02-15 22:20:21', '2022-02-15 22:20:21', '{"display": "display-0010"}', 'force_subject_id-0010', 'auth_session-0010', 'challenge-0010', 'requested_audience-0010_1', 'acr-0010', '{"context": "0010"}'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request_handled
(challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
('challenge-0010', 'granted_scope-0010_1', true, 0010, '{}', '2022-02-15 22:20:21', '{"session_access_token-0010": "0010"}', '{"session_id_token-0010": "0010"}', '2022-02-15 22:20:21', true, 'granted_audience-0010_1');

INSERT INTO hydra_oauth2_authentication_request_handled
(challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier, context)
VALUES
('challenge-0010', 'subject-0010', true, 0010, '{}', 'acr-0010', '2022-02-15 22:20:21', '2022-02-15 22:20:21', true, 'force_subject_id-0010', '{"context": "0010"}');

INSERT INTO hydra_oauth2_obfuscated_authentication_session (client_id, subject, subject_obfuscated)
SELECT hydra_client.id, 'subject-0010', 'subject_obfuscated-0010'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_logout_request (challenge, verifier, subject, sid, client_id, request_url, redir_url, was_used, accepted, rejected, rp_initiated)
SELECT 'challenge-0010', 'verifier-0010', 'subject-0010', 'session_id-0010', hydra_client.id, 'http://request/0010', 'http://post_logout/0010', true, true, false, true
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
