INSERT INTO hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
SELECT 'sig-0003', 'req-0003', now(), hydra_client.id, 'scope-0003', 'granted_scope-0003', 'form_data-0003', 'session-0003', 'subject-0003'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
SELECT 'sig-0003', 'req-0003', now(), hydra_client.id, 'scope-0003', 'granted_scope-0003', 'form_data-0003', 'session-0003', 'subject-0003'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
SELECT 'sig-0003', 'req-0003', now(), hydra_client.id, 'scope-0003', 'granted_scope-0003', 'form_data-0003', 'session-0003', 'subject-0003'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
SELECT 'sig-0003', 'req-0003', now(), hydra_client.id, 'scope-0003', 'granted_scope-0003', 'form_data-0003', 'session-0003', 'subject-0003'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
SELECT 'sig-0003', 'req-0003', now(), hydra_client.id, 'scope-0003', 'granted_scope-0003', 'form_data-0003', 'session-0003', 'subject-0003'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
