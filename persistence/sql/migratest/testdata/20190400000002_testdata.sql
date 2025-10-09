INSERT INTO hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
SELECT 'sig-0002', 'req-0002', '2022-02-15 22:20:21', hydra_client.id, 'scope-0002', 'granted_scope-0002', 'form_data-0002', 'session-0002', 'subject-0002'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
SELECT 'sig-0002', 'req-0002', '2022-02-15 22:20:21', hydra_client.id, 'scope-0002', 'granted_scope-0002', 'form_data-0002', 'session-0002', 'subject-0002'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
SELECT 'sig-0002', 'req-0002', '2022-02-15 22:20:21', hydra_client.id, 'scope-0002', 'granted_scope-0002', 'form_data-0002', 'session-0002', 'subject-0002'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
SELECT 'sig-0002', 'req-0002', '2022-02-15 22:20:21', hydra_client.id, 'scope-0002', 'granted_scope-0002', 'form_data-0002', 'session-0002', 'subject-0002'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
