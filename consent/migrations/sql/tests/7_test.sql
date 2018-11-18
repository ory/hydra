-- +migrate Up
INSERT INTO
	hydra_oauth2_authentication_session (id, authenticated_at, subject)
VALUES
	('7-login-session-id', NOW(), '7-sub');

INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id, requested_at_audience)
VALUES
	('7-challenge', '7-verifier', '7-client', '7-subject', '7-redirect', false, '7-scope', '7-csrf', NOW(), NOW(), '{}', '7-login-session-id', '7-aud');

INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr)
VALUES
	('7-challenge', '7-verifier', '7-client', '7-subject', '7-redirect', false, '7-scope', '7-csrf', NOW(), NOW(), '{}', '7-forced-sub', '7-login-session-id', '7-challenge', '7-aud', '7-acr');

INSERT INTO
	hydra_oauth2_consent_request_handled (challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
	('7-challenge', '7-scope', true, 3600, '{}', NOW(), '{}', '{}', NOW(), false, '7-aud');

INSERT INTO
	hydra_oauth2_authentication_request_handled (challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
	('7-challenge', '7-sub', true, 3600, '{}', '1', NOW(), NOW(), false, '7-forced-sub');

INSERT INTO
	hydra_oauth2_obfuscated_authentication_session (subject, client_id, subject_obfuscated)
VALUES
	('7-sub', '7-client', '7-obfuscated');

-- +migrate Down
