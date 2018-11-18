-- +migrate Up
INSERT INTO
	hydra_oauth2_authentication_session (id, authenticated_at, subject)
VALUES
	('1-login-session-id', NOW(), '1-sub');

INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context)
VALUES
	('1-challenge', '1-verifier', '1-client', '1-subject', '1-redirect', false, '1-scope', '1-csrf', NOW(), NOW(), '{}');

INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context)
VALUES
	('1-challenge', '1-verifier', '1-client', '1-subject', '1-redirect', false, '1-scope', '1-csrf', NOW(), NOW(), '{}');

INSERT INTO
	hydra_oauth2_consent_request_handled (challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used)
VALUES
	('1-challenge', '1-scope', true, 3600, '{}', NOW(), '{}', '{}', NOW(), false);

INSERT INTO
 hydra_oauth2_authentication_request_handled (challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used)
VALUES
	('1-challenge', '1-sub', true, 3600, '{}', '1', NOW(), NOW(), false);

-- +migrate Down
