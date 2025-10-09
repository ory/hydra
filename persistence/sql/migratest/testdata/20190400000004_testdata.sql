INSERT INTO hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
SELECT 'sig-0004', 'req-0004', '2022-02-15 22:20:21', hydra_client.id, 'scope-0004', 'granted_scope-0004', 'form_data-0004', 'session-0004', 'subject-0004', false
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
SELECT 'sig-0004', 'req-0004', '2022-02-15 22:20:21', hydra_client.id, 'scope-0004', 'granted_scope-0004', 'form_data-0004', 'session-0004', 'subject-0004', false
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
SELECT 'sig-0004', 'req-0004', '2022-02-15 22:20:21', hydra_client.id, 'scope-0004', 'granted_scope-0004', 'form_data-0004', 'session-0004', 'subject-0004', false
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
SELECT 'sig-0004', 'req-0004', '2022-02-15 22:20:21', hydra_client.id, 'scope-0004', 'granted_scope-0004', 'form_data-0004', 'session-0004', 'subject-0004', false
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
SELECT 'sig-0004', 'req-0004', '2022-02-15 22:20:21', hydra_client.id, 'scope-0004', 'granted_scope-0004', 'form_data-0004', 'session-0004', 'subject-0004', false
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
