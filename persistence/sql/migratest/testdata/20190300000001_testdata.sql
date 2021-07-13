-- using the most lately added client as a foreign key
INSERT INTO hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context)
SELECT 'challenge-0001', 'verifier-0001', hydra_client.id, 'subject-0001', 'http://request/0001', true, 'requested_scope-0001_1', 'csrf-0001', now(), now(), '{"display": "display-0001"}'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context)
SELECT 'challenge-0001', 'verifier-0001', hydra_client.id, 'subject-0001', 'http://request/0001', true, 'requested_scope-0001_1', 'csrf-0001', now(), now(), '{"display": "display-0001"}'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_authentication_session
(id, authenticated_at, subject)
VALUES
('auth_session-0001', now(), 'subject-0001');

INSERT INTO hydra_oauth2_consent_request_handled
(challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used)
VALUES
('challenge-0001', 'granted_scope-0001_1', true, 0001, '{}', now(), '{"session_access_token-0001": "0001"}', '{"session_id_token-0001": "0001"}', now(), true);

INSERT INTO hydra_oauth2_authentication_request_handled
(challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used)
VALUES
('challenge-0001', 'subject-0001', true, 0001, '{}', 'acr-0001', now(), now(), true);
