-- +migrate Up
INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier)
VALUES
	('2-challenge', '2-verifier', '2-client', '2-subject', '2-redirect', false, '2-scope', '2-csrf', NOW(), NOW(), '{}', '2-forced-sub');

INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context)
VALUES
	('2-challenge', '2-verifier', '2-client', '2-subject', '2-redirect', false, '2-scope', '2-csrf', NOW(), NOW(), '{}');

INSERT INTO
	hydra_oauth2_authentication_session (id, authenticated_at, subject)
VALUES
	('2-auth', NOW(), '2-sub');

INSERT INTO
	hydra_oauth2_consent_request_handled (challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used)
VALUES
	('2-challenge', '2-scope', true, 3600, '{}', NOW(), '{}', '{}', NOW(), false);

INSERT INTO
	hydra_oauth2_authentication_request_handled (challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
	('2-challenge', '2-sub', true, 3600, '{}', '1', NOW(), NOW(), false, '2-forced-sub');

INSERT INTO
	hydra_oauth2_obfuscated_authentication_session (subject, client_id, subject_obfuscated)
VALUES
	('2-sub', '2-client', '2-obfuscated');

-- +migrate Down
