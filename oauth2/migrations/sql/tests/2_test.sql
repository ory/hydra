-- +migrate Up
INSERT INTO
	hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('2-sig', '2-request', NOW(), '2-client', '2-scope', '2-granted-scope', '', '{}', '2-subject');

INSERT INTO
	hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('2-sig', '2-request', NOW(), '2-client', '2-scope', '2-granted-scope', '', '{}', '2-subject');

INSERT INTO
	hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('2-sig', '2-request', NOW(), '2-client', '2-scope', '2-granted-scope', '', '{}', '2-subject');

INSERT INTO
	hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('2-sig', '2-request', NOW(), '2-client', '2-scope', '2-granted-scope', '', '{}', '2-subject');

-- +migrate Down
