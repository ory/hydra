INSERT INTO hydra_oauth2_authentication_session
(id, authenticated_at, subject, remember)
VALUES
('auth_session-0011', now(), 'subject-0011', false);

-- using the most lately added client as a foreign key
INSERT INTO hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, requested_at_audience, login_session_id)
SELECT 'challenge-0011', 'verifier-0011', hydra_client.id, 'subject-0011', 'http://request/0011', true, 'requested_scope-0011_1', 'csrf-0011', now(), now(), '{"display": "display-0011"}', 'requested_audience-0011_1', 'auth_session-0011'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr, context)
SELECT 'challenge-0011', 'verifier-0011', hydra_client.id, 'subject-0011', 'http://request/0011', true, 'requested_scope-0011_1', 'csrf-0011', now(), now(), '{"display": "display-0011"}', 'force_subject_id-0011', 'auth_session-0011', 'challenge-0011', 'requested_audience-0011_1', 'acr-0011', '{"context": "0011"}'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request_handled
(challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
('challenge-0011', 'granted_scope-0011_1', true, 0011, '{}', now(), '{"session_access_token-0011": "0011"}', '{"session_id_token-0011": "0011"}', now(), true, 'granted_audience-0011_1');

INSERT INTO hydra_oauth2_authentication_request_handled
(challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier, context)
VALUES
('challenge-0011', 'subject-0011', true, 0011, '{}', 'acr-0011', now(), now(), true, 'force_subject_id-0011', '{"context": "0011"}');

INSERT INTO hydra_oauth2_obfuscated_authentication_session (client_id, subject, subject_obfuscated)
SELECT hydra_client.id, 'subject-0011', 'subject_obfuscated-0011'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_logout_request (challenge, verifier, subject, sid, client_id, request_url, redir_url, was_used, accepted, rejected, rp_initiated)
SELECT 'challenge-0011', 'verifier-0011', 'subject-0011', 'session_id-0011', hydra_client.id, 'http://request/0011', 'http://post_logout/0011', true, true, false, true
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
