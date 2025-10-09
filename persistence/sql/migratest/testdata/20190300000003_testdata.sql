-- 20190300000001_testdata.sql (see 20190300000001_testdata.sql for an explanation)
-- using the most lately added client as a foreign key
INSERT INTO hydra_oauth2_consent_request (challenge, login_challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context)
SELECT 'challenge-0001', 'challenge-0001', 'verifier-0001', hydra_client.id, 'subject-0001', 'http://request/0001', true, 'requested_scope-0001_1', 'csrf-0001', '2022-02-15 22:20:21', '2022-02-15 22:20:21', '{"display": "display-0001"}'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context)
SELECT 'challenge-0001', 'verifier-0001', hydra_client.id, 'subject-0001', 'http://request/0001', true, 'requested_scope-0001_1', 'csrf-0001', '2022-02-15 22:20:21', '2022-02-15 22:20:21', '{"display": "display-0001"}'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_authentication_session
(id, authenticated_at, subject)
VALUES
('auth_session-0001', '2022-02-15 22:20:21', 'subject-0001');

INSERT INTO hydra_oauth2_consent_request_handled
(challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used)
VALUES
('challenge-0001', 'granted_scope-0001_1', true, 0001, '{}', '2022-02-15 22:20:21', '{"session_access_token-0001": "0001"}', '{"session_id_token-0001": "0001"}', '2022-02-15 22:20:21', true);

INSERT INTO hydra_oauth2_authentication_request_handled
(challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used)
VALUES
('challenge-0001', 'subject-0001', true, 0001, '{}', 'acr-0001', '2022-02-15 22:20:21', '2022-02-15 22:20:21', true);
-- EOF 20190300000001_testdata.sql

-- 20190300000002_testdata.sql (see 20190300000002_testdata.sql for an explanation)
-- using the most lately added client as a foreign key
INSERT INTO hydra_oauth2_consent_request (challenge, login_challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier)
SELECT 'challenge-0002', 'challenge-0002', 'verifier-0002', hydra_client.id, 'subject-0002', 'http://request/0002', true, 'requested_scope-0002_1', 'csrf-0002', '2022-02-15 22:20:21', '2022-02-15 22:20:21', '{"display": "display-0002"}', 'force_subject_id-0002'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context)
SELECT 'challenge-0002', 'verifier-0002', hydra_client.id, 'subject-0002', 'http://request/0002', true, 'requested_scope-0002_1', 'csrf-0002', '2022-02-15 22:20:21', '2022-02-15 22:20:21', '{"display": "display-0002"}'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_authentication_session
(id, authenticated_at, subject)
VALUES
('auth_session-0002', '2022-02-15 22:20:21', 'subject-0002');

INSERT INTO hydra_oauth2_consent_request_handled
(challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used)
VALUES
('challenge-0002', 'granted_scope-0002_1', true, 0002, '{}', '2022-02-15 22:20:21', '{"session_access_token-0002": "0002"}', '{"session_id_token-0002": "0002"}', '2022-02-15 22:20:21', true);

INSERT INTO hydra_oauth2_authentication_request_handled
(challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
('challenge-0002', 'subject-0002', true, 0002, '{}', 'acr-0002', '2022-02-15 22:20:21', '2022-02-15 22:20:21', true, 'force_subject_id-0002');

INSERT INTO hydra_oauth2_obfuscated_authentication_session (client_id, subject, subject_obfuscated)
SELECT hydra_client.id, 'subject-0002', 'subject_obfuscated-0002'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
-- EOF 20190300000002_testdata.sql


INSERT INTO hydra_oauth2_authentication_session
(id, authenticated_at, subject)
VALUES
('auth_session-0003', '2022-02-15 22:20:21', 'subject-0003');

-- using the most lately added client as a foreign key
INSERT INTO hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id)
SELECT 'challenge-0003', 'verifier-0003', hydra_client.id, 'subject-0003', 'http://request/0003', true, 'requested_scope-0003_1', 'csrf-0003', '2022-02-15 22:20:21', '2022-02-15 22:20:21', '{"display": "display-0003"}', 'auth_session-0003'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge)
SELECT 'challenge-0003', 'verifier-0003', hydra_client.id, 'subject-0003', 'http://request/0003', true, 'requested_scope-0003_1', 'csrf-0003', '2022-02-15 22:20:21', '2022-02-15 22:20:21', '{"display": "display-0003"}', 'force_subject_id-0003', 'auth_session-0003', 'challenge-0003'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_consent_request_handled
(challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used)
VALUES
('challenge-0003', 'granted_scope-0003_1', true, 0003, '{}', '2022-02-15 22:20:21', '{"session_access_token-0003": "0003"}', '{"session_id_token-0003": "0003"}', '2022-02-15 22:20:21', true);

INSERT INTO hydra_oauth2_authentication_request_handled
(challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
('challenge-0003', 'subject-0003', true, 0003, '{}', 'acr-0003', '2022-02-15 22:20:21', '2022-02-15 22:20:21', true, 'force_subject_id-0003');

INSERT INTO hydra_oauth2_obfuscated_authentication_session (client_id, subject, subject_obfuscated)
SELECT hydra_client.id, 'subject-0003', 'subject_obfuscated-0003'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
