-- +migrate Up
INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience)
VALUES
	('5-challenge', '5-verifier', '5-client', '5-subject', '5-redirect', false, '5-scope', '5-csrf', NOW(), NOW(), '{}', '5-forced-sub', '5-login_session_id', '5-login-challenge', '5-aud');

INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id, requested_at_audience)
VALUES
	('5-challenge', '5-verifier', '5-client', '5-subject', '5-redirect', false, '5-scope', '5-csrf', NOW(), NOW(), '{}', '5-login-session-id', '5-aud');

INSERT INTO
	hydra_oauth2_authentication_session (id, authenticated_at, subject)
VALUES
	('5-auth', NOW(), '5-sub');

INSERT INTO
	hydra_oauth2_consent_request_handled (challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
	('5-challenge', '5-scope', true, 3600, '{}', NOW(), '{}', '{}', NOW(), false, '5-aud');

INSERT INTO
	hydra_oauth2_authentication_request_handled (challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
	('5-challenge', '5-sub', true, 3600, '{}', '1', NOW(), NOW(), false, '5-forced-sub');

INSERT INTO
	hydra_oauth2_obfuscated_authentication_session (subject, client_id, subject_obfuscated)
VALUES
	('5-sub', '5-client', '5-obfuscated');

-- +migrate Down
