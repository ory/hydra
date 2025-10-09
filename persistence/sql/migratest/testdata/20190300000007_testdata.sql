INSERT INTO hydra_oauth2_authentication_session
(id, authenticated_at, subject)
VALUES
('auth_session-0007', '2022-02-15 22:20:21', 'subject-0007');

-- using the most lately added client as a foreign key
INSERT INTO hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, requested_at_audience, login_session_id)
SELECT 'challenge-0007', 'verifier-0007', hydra_client.id, 'subject-0007', 'http://request/0007', true, 'requested_scope-0007_1', 'csrf-0007', '2022-02-15 22:20:21', '2022-02-15 22:20:21', '{"display": "display-0007"}', 'requested_audience-0007_1', 'auth_session-0007'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr)
SELECT 'challenge-0007', 'verifier-0007', hydra_client.id, 'subject-0007', 'http://request/0007', true, 'requested_scope-0007_1', 'csrf-0007', '2022-02-15 22:20:21', '2022-02-15 22:20:21', '{"display": "display-0007"}', 'force_subject_id-0007', 'auth_session-0007', 'challenge-0007', 'requested_audience-0007_1', 'acr-0007'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request_handled
(challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
('challenge-0007', 'granted_scope-0007_1', true, 0007, '{}', '2022-02-15 22:20:21', '{"session_access_token-0007": "0007"}', '{"session_id_token-0007": "0007"}', '2022-02-15 22:20:21', true, 'granted_audience-0007_1');

INSERT INTO hydra_oauth2_authentication_request_handled
(challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
('challenge-0007', 'subject-0007', true, 0007, '{}', 'acr-0007', '2022-02-15 22:20:21', '2022-02-15 22:20:21', true, 'force_subject_id-0007');

INSERT INTO hydra_oauth2_obfuscated_authentication_session (client_id, subject, subject_obfuscated)
SELECT hydra_client.id, 'subject-0007', 'subject_obfuscated-0007'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
