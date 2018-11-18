-- +migrate Up
INSERT INTO hydra_client (id, allowed_cors_origins, client_name, client_secret, redirect_uris, grant_types, response_types, scope, owner, policy_uri, tos_uri, client_uri, logo_uri, contacts, client_secret_expires_at, sector_identifier_uri, jwks, jwks_uri, token_endpoint_auth_method, request_uris, request_object_signing_alg, userinfo_signed_response_alg, subject_type, audience)
VALUES
  ('6-client', 'http://localhost|http://google', 'some-client', 'abcdef', 'http://localhost|http://google', 'authorize_code|implicit', 'token|id_token', 'foo|bar', 'aeneas', 'http://policy', 'http://tos', 'http://client', 'http://logo', 'aeneas|foo', 0, 'http://sector', '{"keys": []}', 'http://jwks', 'none', 'http://uri1|http://uri2', 'rs256', 'rs526', 'public', 'https://www.ory.sh/api');

INSERT INTO
	hydra_oauth2_authentication_session (id, authenticated_at, subject)
VALUES
	('6-login-session-id', NOW(), '6-sub');

INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id, requested_at_audience)
VALUES
	('6-challenge', '6-verifier', '6-client', '6-subject', '6-redirect', false, '6-scope', '6-csrf', NOW(), NOW(), '{}', '6-login-session-id', '6-aud');

INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr)
VALUES
	('6-challenge', '6-verifier', '6-client', '6-subject', '6-redirect', false, '6-scope', '6-csrf', NOW(), NOW(), '{}', '6-forced-sub', '6-login-session-id', '6-challenge', '6-aud', '6-acr');

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

-- This is inconsistent data specifically inserted for the case where we have inconsistent state before executing migration 7 which adds foreign key constraints.
INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id, requested_at_audience)
VALUES
	('6-bfk1-challenge', '6-bfk1-verifier', '6-client', '6-subject', '6-redirect', false, '6-scope', '6-csrf', NOW(), NOW(), '{}', 'i-do-not-exist', '6-aud');
INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id, requested_at_audience)
VALUES
	('6-bfk2-challenge', '6-bfk2-verifier', 'i-do-not-exist', '6-subject', '6-redirect', false, '6-scope', '6-csrf', NOW(), NOW(), '{}', '6-login-session-id', '6-aud');

INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr)
VALUES
	('6-bfk1-challenge', '6-bfk1-verifier', '6-client', '6-subject', '6-redirect', false, '6-scope', '6-csrf', NOW(), NOW(), '{}', '6-forced-sub', '6-login-session-id', 'i-do-not-exist', '6-aud', '6-acr');
INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr)
VALUES
	('6-bfk2-challenge', '6-bfk2-verifier', '6-client', '6-subject', '6-redirect', false, '6-scope', '6-csrf', NOW(), NOW(), '{}', '6-forced-sub', 'i-do-not-exist', '6-challenge', '6-aud', '6-acr');
INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr)
VALUES
	('6-bfk3-challenge', '6-bfk3-verifier', 'i-do-not-exist', '6-subject', '6-redirect', false, '6-scope', '6-csrf', NOW(), NOW(), '{}', '6-forced-sub', '6-login-session-id', '6-challenge', '6-aud', '6-acr');

INSERT INTO
	hydra_oauth2_consent_request_handled (challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
	('i-do-not-exist', '6-bfk-scope', true, 3600, '{}', NOW(), '{}', '{}', NOW(), false, '6-bfk-aud');

INSERT INTO
	hydra_oauth2_authentication_request_handled (challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier)
VALUES
	('i-do-not-exist', '6-bfk-sub', true, 3600, '{}', '1', NOW(), NOW(), false, '6-bfk-forced-sub');

INSERT INTO
	hydra_oauth2_obfuscated_authentication_session (subject, client_id, subject_obfuscated)
VALUES
	('6-bfk-sub', 'i-do-not-exist', '6-bfk-obfuscated');


-- +migrate Down
