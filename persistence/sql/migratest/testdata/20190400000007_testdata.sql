INSERT INTO hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience)
SELECT 'sig-0007', 'req-0007', '2022-02-15 22:20:22', hydra_client.id, 'scope-0007', 'granted_scope-0007', 'form_data-0007', 'session-0007', 'subject-0007', false, 'requested_audience-0007', 'granted_audience-0007'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience)
SELECT 'sig-0007', 'req-0007', '2022-02-15 22:20:22', hydra_client.id, 'scope-0007', 'granted_scope-0007', 'form_data-0007', 'session-0007', 'subject-0007', false, 'requested_audience-0007', 'granted_audience-0007'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience)
SELECT 'sig-0007', 'req-0007', '2022-02-15 22:20:22', hydra_client.id, 'scope-0007', 'granted_scope-0007', 'form_data-0007', 'session-0007', 'subject-0007', false, 'requested_audience-0007', 'granted_audience-0007'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience)
SELECT 'sig-0007', 'req-0007', '2022-02-15 22:20:22', hydra_client.id, 'scope-0007', 'granted_scope-0007', 'form_data-0007', 'session-0007', 'subject-0007', false, 'requested_audience-0007', 'granted_audience-0007'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;

INSERT INTO hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience)
SELECT 'sig-0007', 'req-0007', '2022-02-15 22:20:22', hydra_client.id, 'scope-0007', 'granted_scope-0007', 'form_data-0007', 'session-0007', 'subject-0007', false, 'requested_audience-0007', 'granted_audience-0007'
FROM hydra_client
ORDER BY hydra_client.pk DESC
LIMIT 1;
