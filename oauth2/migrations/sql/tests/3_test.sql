-- +migrate Up
INSERT INTO hydra_client (id, allowed_cors_origins, client_name, client_secret, redirect_uris, grant_types, response_types, scope, owner, policy_uri, tos_uri, client_uri, logo_uri, contacts, client_secret_expires_at, sector_identifier_uri, jwks, jwks_uri, token_endpoint_auth_method, request_uris, request_object_signing_alg, userinfo_signed_response_alg, subject_type, audience)
VALUES
  ('3-client', 'http://localhost|http://google', 'some-client', 'abcdef', 'http://localhost|http://google', 'authorize_code|implicit', 'token|id_token', 'foo|bar', 'aeneas', 'http://policy', 'http://tos', 'http://client', 'http://logo', 'aeneas|foo', 0, 'http://sector', '{"keys": []}', 'http://jwks', 'none', 'http://uri1|http://uri2', 'rs256', 'rs526', 'public', 'https://www.ory.sh/api');

INSERT INTO
	hydra_oauth2_authentication_session (id, authenticated_at, subject)
VALUES
	('3-login-session-id', NOW(), '3-sub');

INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id, requested_at_audience)
VALUES
	('3-request', '3-verifier', '3-client', '3-subject', '3-redirect', false, '3-scope', '3-csrf', NOW(), NOW(), '{}', '3-login-session-id', '3-aud');

INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr)
VALUES
	('3-request', '3-verifier', '3-client', '3-subject', '3-redirect', false, '3-scope', '3-csrf', NOW(), NOW(), '{}', '3-forced-sub', '3-login-session-id', '3-request', '3-aud', '3-acr');

INSERT INTO
	hydra_oauth2_consent_request_handled (challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
	('3-request', '3-scope', true, 3600, '{}', NOW(), '{}', '{}', NOW(), false, '3-aud');

-- The previous block is just to get foreign keys working

INSERT INTO
	hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('3-sig', '3-request', NOW(), '3-client', '3-scope', '3-granted-scope', '', '{}', '3-subject');

INSERT INTO
	hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('3-sig', '3-request', NOW(), '3-client', '3-scope', '3-granted-scope', '', '{}', '3-subject');

INSERT INTO
	hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('3-sig', '3-request', NOW(), '3-client', '3-scope', '3-granted-scope', '', '{}', '3-subject');

INSERT INTO
	hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('3-sig', '3-request', NOW(), '3-client', '3-scope', '3-granted-scope', '', '{}', '3-subject');

INSERT INTO
	hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject)
VALUES
	('3-sig', '3-request', NOW(), '3-client', '3-scope', '3-granted-scope', '', '{}', '3-subject');

-- +migrate Down
