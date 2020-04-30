-- +migrate Up
INSERT INTO
	hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data)
VALUES
	('1-sig', '1-request', NOW(), '1-client', '1-scope', '1-granted-scope', '', '{}');

INSERT INTO
	hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data)
VALUES
	('1-sig', '1-request', NOW(), '1-client', '1-scope', '1-granted-scope', '', '{}');

INSERT INTO
	hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data)
VALUES
	('1-sig', '1-request', NOW(), '1-client', '1-scope', '1-granted-scope', '', '{}');

INSERT INTO
	hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data)
VALUES
	('1-sig', '1-request', NOW(), '1-client', '1-scope', '1-granted-scope', '', '{}');

-- +migrate Down
