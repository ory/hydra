-- +migrate Up
INSERT INTO hydra_client (id, allowed_cors_origins, client_name, client_secret, redirect_uris, grant_types, response_types, scope, owner, policy_uri, tos_uri, client_uri, logo_uri, contacts, client_secret_expires_at, sector_identifier_uri, jwks, jwks_uri, token_endpoint_auth_method, request_uris, request_object_signing_alg, userinfo_signed_response_alg, subject_type, audience)
VALUES
  ('9-client', 'http://localhost|http://google', 'some-client', 'abcdef', 'http://localhost|http://google', 'authorize_code|implicit', 'token|id_token', 'foo|bar', 'aeneas', 'http://policy', 'http://tos', 'http://client', 'http://logo', 'aeneas|foo', 0, 'http://sector', '{"keys": []}', 'http://jwks', 'none', 'http://uri1|http://uri2', 'rs256', 'rs526', 'public', 'https://www.ory.sh/api');

INSERT INTO
	hydra_oauth2_authentication_session (id, authenticated_at, subject)
VALUES
	('9-login-session-id', NOW(), '9-sub');

INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id, requested_at_audience)
VALUES
	('9-challenge', '9-verifier', '9-client', '9-subject', '9-redirect', false, '9-scope', '9-csrf', NOW(), NOW(), '{}', '9-login-session-id', '9-aud');

INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr)
VALUES
	('9-challenge', '9-verifier', '9-client', '9-subject', '9-redirect', false, '9-scope', '9-csrf', NOW(), NOW(), '{}', '9-forced-sub', '9-login-session-id', '9-challenge', '9-aud', '9-acr');

INSERT INTO
	hydra_oauth2_consent_request_handled (challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
	('9-challenge', '9-scope', true, 3600, '{}', NOW(), '{}', '{}', NOW(), false, '9-aud');

-- The previous block is just to get foreign keys working

INSERT INTO
	hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
VALUES
	('9-sig', '9-request', NOW(), '9-client', '9-scope', '9-granted-scope', '', '{}', '9-subject', true, '9-challengeed-aud', '9-granted-aud', '9-challenge');

INSERT INTO
	hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
VALUES
	('9-sig', '9-request', NOW(), '9-client', '9-scope', '9-granted-scope', '', '{}', '9-subject', true, '9-challengeed-aud', '9-granted-aud', '9-challenge');

INSERT INTO
	hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
VALUES
	('9-sig', '9-request', NOW(), '9-client', '9-scope', '9-granted-scope', '', '{}', '9-subject', true, '9-challengeed-aud', '9-granted-aud', '9-challenge');

INSERT INTO
	hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
VALUES
	('9-sig', '9-request', NOW(), '9-client', '9-scope', '9-granted-scope', '', '{}', '9-subject', true, '9-challengeed-aud', '9-granted-aud', '9-challenge');

INSERT INTO
	hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
VALUES
	('9-sig', '9-request', NOW(), '9-client', '9-scope', '9-granted-scope', '', '{}', '9-subject', true, '9-challengeed-aud', '9-granted-aud', '9-challenge');

-- +migrate Down
