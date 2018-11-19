-- +migrate Up

INSERT INTO
	hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience)
VALUES
	('7-sig', '7-request', NOW(), '7-client', '7-scope', '7-granted-scope', '', '{}', '7-subject', true, '7-requested-aud', '7-granted-aud');

INSERT INTO
	hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience)
VALUES
	('7-sig', '7-request', NOW(), '7-client', '7-scope', '7-granted-scope', '', '{}', '7-subject', true, '7-requested-aud', '7-granted-aud');

INSERT INTO
	hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience)
VALUES
	('7-sig', '7-request', NOW(), '7-client', '7-scope', '7-granted-scope', '', '{}', '7-subject', true, '7-requested-aud', '7-granted-aud');

INSERT INTO
	hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience)
VALUES
	('7-sig', '7-request', NOW(), '7-client', '7-scope', '7-granted-scope', '', '{}', '7-subject', true, '7-requested-aud', '7-granted-aud');

INSERT INTO
	hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience)
VALUES
	('7-sig', '7-request', NOW(), '7-client', '7-scope', '7-granted-scope', '', '{}', '7-subject', true, '7-requested-aud', '7-granted-aud');

-- +migrate Down
