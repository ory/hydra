INSERT INTO hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
SELECT 'sig-0005', 'req-0005', '2022-02-15 22:20:22', hydra_client.id, 'scope-0005', 'granted_scope-0005', 'form_data-0005', 'session-0005', 'subject-0005', false
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
SELECT 'sig-0005', 'req-0005', '2022-02-15 22:20:22', hydra_client.id, 'scope-0005', 'granted_scope-0005', 'form_data-0005', 'session-0005', 'subject-0005', false
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
SELECT 'sig-0005', 'req-0005', '2022-02-15 22:20:22', hydra_client.id, 'scope-0005', 'granted_scope-0005', 'form_data-0005', 'session-0005', 'subject-0005', false
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
SELECT 'sig-0005', 'req-0005', '2022-02-15 22:20:22', hydra_client.id, 'scope-0005', 'granted_scope-0005', 'form_data-0005', 'session-0005', 'subject-0005', false
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
SELECT 'sig-0005', 'req-0005', '2022-02-15 22:20:22', hydra_client.id, 'scope-0005', 'granted_scope-0005', 'form_data-0005', 'session-0005', 'subject-0005', false
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
