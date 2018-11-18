-- +migrate Up
INSERT INTO
	hydra_oauth2_authentication_session (id, authenticated_at, subject)
VALUES
	('4-login-session-id', NOW(), '4-sub');

INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id, requested_at_audience)
VALUES
	('4-challenge', '4-verifier', '4-client', '4-subject', '4-redirect', false, '4-scope', '4-csrf', NOW(), NOW(), '{}', '4-login-session-id', '4-aud');

INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience)
VALUES
	('4-challenge', '4-verifier', '4-client', '4-subject', '4-redirect', false, '4-scope', '4-csrf', NOW(), NOW(), '{}', '4-forced-sub', '4-login-session-id', '4-challenge', '4-aud');

INSERT INTO
	hydra_oauth2_consent_request_handled (challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
	('4-challenge', '4-scope', true, 3600, '{}', NOW(), '{}', '{}', NOW(), false, '4-aud');

INSERT INTO
	hydra_oauth2_authentication_request_handled (challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
	('4-challenge', '4-sub', true, 3600, '{}', '1', NOW(), NOW(), false, '4-forced-sub');

INSERT INTO
	hydra_oauth2_obfuscated_authentication_session (subject, client_id, subject_obfuscated)
VALUES
	('4-sub', '4-client', '4-obfuscated');

-- +migrate Down
