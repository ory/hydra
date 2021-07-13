-- using the most lately added client as a foreign key
INSERT INTO hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier)
SELECT 'challenge-0002', 'verifier-0002', hydra_client.id, 'subject-0002', 'http://request/0002', true, 'requested_scope-0002_1', 'csrf-0002', now(), now(), '{"display": "display-0002"}', 'force_subject_id-0002'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context)
SELECT 'challenge-0002', 'verifier-0002', hydra_client.id, 'subject-0002', 'http://request/0002', true, 'requested_scope-0002_1', 'csrf-0002', now(), now(), '{"display": "display-0002"}'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_authentication_session
(id, authenticated_at, subject)
VALUES
('auth_session-0002', now(), 'subject-0002');

INSERT INTO hydra_oauth2_consent_request_handled
(challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used)
VALUES
('challenge-0002', 'granted_scope-0002_1', true, 0002, '{}', now(), '{"session_access_token-0002": "0002"}', '{"session_id_token-0002": "0002"}', now(), true);

INSERT INTO hydra_oauth2_authentication_request_handled
(challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
('challenge-0002', 'subject-0002', true, 0002, '{}', 'acr-0002', now(), now(), true, 'force_subject_id-0002');

INSERT INTO hydra_oauth2_obfuscated_authentication_session (client_id, subject, subject_obfuscated)
SELECT hydra_client.id, 'subject-0002', 'subject_obfuscated-0002'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
