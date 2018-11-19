-- +migrate Up

INSERT INTO
	hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('6-sig', '6-request', NOW(), '6-client', '6-scope', '6-granted-scope', '', '{}', '6-subject', true);

INSERT INTO
	hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('6-sig', '6-request', NOW(), '6-client', '6-scope', '6-granted-scope', '', '{}', '6-subject', true);

INSERT INTO
	hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('6-sig', '6-request', NOW(), '6-client', '6-scope', '6-granted-scope', '', '{}', '6-subject', true);

INSERT INTO
	hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('6-sig', '6-request', NOW(), '6-client', '6-scope', '6-granted-scope', '', '{}', '6-subject', true);

INSERT INTO
	hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('6-sig', '6-request', NOW(), '6-client', '6-scope', '6-granted-scope', '', '{}', '6-subject', true);

-- +migrate Down
