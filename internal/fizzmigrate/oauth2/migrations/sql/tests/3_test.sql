-- +migrate Up

INSERT INTO
	hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('3-sig', '3-request', NOW(), '3-client', '3-scope', '3-granted-scope', '', '{}', '3-subject');

INSERT INTO
	hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('3-sig', '3-request', NOW(), '3-client', '3-scope', '3-granted-scope', '', '{}', '3-subject');

INSERT INTO
	hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('3-sig', '3-request', NOW(), '3-client', '3-scope', '3-granted-scope', '', '{}', '3-subject');

INSERT INTO
	hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('3-sig', '3-request', NOW(), '3-client', '3-scope', '3-granted-scope', '', '{}', '3-subject');

INSERT INTO
	hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('3-sig', '3-request', NOW(), '3-client', '3-scope', '3-granted-scope', '', '{}', '3-subject');

-- +migrate Down
