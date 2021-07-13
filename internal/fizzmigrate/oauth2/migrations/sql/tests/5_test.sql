-- +migrate Up

INSERT INTO
	hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('5-sig', '5-request', NOW(), '5-client', '5-scope', '5-granted-scope', '', '{}', '5-subject', true);

INSERT INTO
	hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('5-sig', '5-request', NOW(), '5-client', '5-scope', '5-granted-scope', '', '{}', '5-subject', true);

INSERT INTO
	hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('5-sig', '5-request', NOW(), '5-client', '5-scope', '5-granted-scope', '', '{}', '5-subject', true);

INSERT INTO
	hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('5-sig', '5-request', NOW(), '5-client', '5-scope', '5-granted-scope', '', '{}', '5-subject', true);

INSERT INTO
	hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('5-sig', '5-request', NOW(), '5-client', '5-scope', '5-granted-scope', '', '{}', '5-subject', true);

-- +migrate Down
