INSERT INTO hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data)
SELECT 'sig-0001', 'req-0001', now(), hydra_client.id, 'scope-0001', 'granted_scope-0001', 'form_data-0001', 'session-0001'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data)
SELECT 'sig-0001', 'req-0001', now(), hydra_client.id, 'scope-0001', 'granted_scope-0001', 'form_data-0001', 'session-0001'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data)
SELECT 'sig-0001', 'req-0001', now(), hydra_client.id, 'scope-0001', 'granted_scope-0001', 'form_data-0001', 'session-0001'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data)
SELECT 'sig-0001', 'req-0001', now(), hydra_client.id, 'scope-0001', 'granted_scope-0001', 'form_data-0001', 'session-0001'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
