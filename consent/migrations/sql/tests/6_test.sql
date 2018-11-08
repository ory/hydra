-- +migrate Up
INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr)
VALUES
	('6-challenge', '6-verifier', '6-client', '6-subject', '6-redirect', false, '6-scope', '6-csrf', NOW(), NOW(), '{}', '6-forced-sub', '6-login_session_id', '6-login-challenge', '6-aud', '6-acr');

INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id, requested_at_audience)
VALUES
	('6-challenge', '6-verifier', '6-client', '6-subject', '6-redirect', false, '6-scope', '6-csrf', NOW(), NOW(), '{}', '6-login-session-id', '6-aud');

INSERT INTO
	hydra_oauth2_authentication_session (id, authenticated_at, subject)
VALUES
	('6-auth', NOW(), '6-sub');

INSERT INTO
	hydra_oauth2_consent_request_handled (challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
	('6-challenge', '6-scope', true, 3600, '{}', NOW(), '{}', '{}', NOW(), false, '6-aud');

INSERT INTO
	hydra_oauth2_authentication_request_handled (challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
	('6-challenge', '6-sub', true, 3600, '{}', '1', NOW(), NOW(), false, '6-forced-sub');

INSERT INTO
	hydra_oauth2_obfuscated_authentication_session (subject, client_id, subject_obfuscated)
VALUES
	('6-sub', '6-client', '6-obfuscated');

-- +migrate Down
