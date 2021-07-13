INSERT INTO hydra_oauth2_authentication_session
(id, authenticated_at, subject, remember)
VALUES
('auth_session-0012', now(), 'subject-0012', false);

-- using the most lately added client as a foreign key
INSERT INTO hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, requested_at_audience, login_session_id)
SELECT 'challenge-0012', 'verifier-0012', hydra_client.id, 'subject-0012', 'http://request/0012', true, 'requested_scope-0012_1', 'csrf-0012', now(), now(), '{"display": "display-0012"}', 'requested_audience-0012_1', 'auth_session-0012'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr, context)
SELECT 'challenge-0012', 'verifier-0012', hydra_client.id, 'subject-0012', 'http://request/0012', true, 'requested_scope-0012_1', 'csrf-0012', now(), now(), '{"display": "display-0012"}', 'force_subject_id-0012', 'auth_session-0012', 'challenge-0012', 'requested_audience-0012_1', 'acr-0012', '{"context": "0012"}'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request_handled
(challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
('challenge-0012', 'granted_scope-0012_1', true, 0012, '{}', now(), '{"session_access_token-0012": "0012"}', '{"session_id_token-0012": "0012"}', now(), true, 'granted_audience-0012_1');

INSERT INTO hydra_oauth2_authentication_request_handled
(challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier, context)
VALUES
('challenge-0012', 'subject-0012', true, 0012, '{}', 'acr-0012', now(), now(), true, 'force_subject_id-0012', '{"context": "0012"}');

INSERT INTO hydra_oauth2_obfuscated_authentication_session (client_id, subject, subject_obfuscated)
SELECT hydra_client.id, 'subject-0012', 'subject_obfuscated-0012'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_logout_request (challenge, verifier, subject, sid, client_id, request_url, redir_url, was_used, accepted, rejected, rp_initiated)
SELECT 'challenge-0012', 'verifier-0012', 'subject-0012', 'session_id-0012', hydra_client.id, 'http://request/0012', 'http://post_logout/0012', true, true, false, true
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
