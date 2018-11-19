-- +migrate Up

INSERT INTO
	hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
VALUES
	('8-sig', '8-request', NOW(), '8-client', '8-scope', '8-granted-scope', '', '{}', '8-subject', true, '8-requested-aud', '8-granted-aud', NULL);

INSERT INTO
	hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
VALUES
	('8-sig', '8-request', NOW(), '8-client', '8-scope', '8-granted-scope', '', '{}', '8-subject', true, '8-requested-aud', '8-granted-aud', NULL);

INSERT INTO
	hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
VALUES
	('8-sig', '8-request', NOW(), '8-client', '8-scope', '8-granted-scope', '', '{}', '8-subject', true, '8-requested-aud', '8-granted-aud', NULL);

INSERT INTO
	hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
VALUES
	('8-sig', '8-request', NOW(), '8-client', '8-scope', '8-granted-scope', '', '{}', '8-subject', true, '8-requested-aud', '8-granted-aud', NULL);

INSERT INTO
	hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
VALUES
	('8-sig', '8-request', NOW(), '8-client', '8-scope', '8-granted-scope', '', '{}', '8-subject', true, '8-requested-aud', '8-granted-aud', NULL);

-- +migrate Down
