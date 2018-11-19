-- +migrate Up

INSERT INTO
	hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('4-sig', '4-request', NOW(), '4-client', '4-scope', '4-granted-scope', '', '{}', '4-subject', true);

INSERT INTO
	hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('4-sig', '4-request', NOW(), '4-client', '4-scope', '4-granted-scope', '', '{}', '4-subject', true);

INSERT INTO
	hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('4-sig', '4-request', NOW(), '4-client', '4-scope', '4-granted-scope', '', '{}', '4-subject', true);

INSERT INTO
	hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('4-sig', '4-request', NOW(), '4-client', '4-scope', '4-granted-scope', '', '{}', '4-subject', true);

INSERT INTO
	hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active)
VALUES
	('4-sig', '4-request', NOW(), '4-client', '4-scope', '4-granted-scope', '', '{}', '4-subject', true);

-- +migrate Down
