-- +migrate Up
INSERT INTO
	hydra_oauth2_authentication_session (id, authenticated_at, subject)
VALUES
	('3-login-session-id', NOW(), '3-sub');

INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id)
VALUES
	('3-challenge', '3-verifier', '3-client', '3-subject', '3-redirect', false, '3-scope', '3-csrf', NOW(), NOW(), '{}', '3-login-session-id');

INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge)
VALUES
	('3-challenge', '3-verifier', '3-client', '3-subject', '3-redirect', false, '3-scope', '3-csrf', NOW(), NOW(), '{}', '3-forced-sub', '3-login-session-id', '3-challenge');

INSERT INTO
	hydra_oauth2_consent_request_handled (challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used)
VALUES
	('3-challenge', '3-scope', true, 3600, '{}', NOW(), '{}', '{}', NOW(), false);

INSERT INTO
	hydra_oauth2_authentication_request_handled (challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
	('3-challenge', '3-sub', true, 3600, '{}', '1', NOW(), NOW(), false, '3-forced-sub');

INSERT INTO
	hydra_oauth2_obfuscated_authentication_session (subject, client_id, subject_obfuscated)
VALUES
	('3-sub', '3-client', '3-obfuscated');

-- +migrate Down
