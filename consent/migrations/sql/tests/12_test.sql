-- +migrate Up
INSERT INTO hydra_client (id, allowed_cors_origins, client_name, client_secret, redirect_uris, grant_types, response_types, scope, owner, policy_uri, tos_uri, client_uri, logo_uri, contacts, client_secret_expires_at, sector_identifier_uri, jwks, jwks_uri, token_endpoint_auth_method, request_uris, request_object_signing_alg, userinfo_signed_response_alg, subject_type, audience, frontchannel_logout_uri, frontchannel_logout_session_required, post_logout_redirect_uris, backchannel_logout_uri, backchannel_logout_session_required)
VALUES
  ('12-client', 'http://localhost|http://google', 'some-client', 'abcdef', 'http://localhost|http://google', 'authorize_code|implicit', 'token|id_token', 'foo|bar', 'aeneas', 'http://policy', 'http://tos', 'http://client', 'http://logo', 'aeneas|foo', 0, 'http://sector', '{"keys": []}', 'http://jwks', 'none', 'http://uri1|http://uri2', 'rs256', 'rs526', 'public', 'https://www.ory.sh/api', 'http://fc-logout/', true, 'http://redir1/|http://redir2/', 'http://bc-logout/', true);

INSERT INTO
	hydra_oauth2_authentication_session (id, authenticated_at, subject, remember)
VALUES
	('12-login-session-id', NOW(), '12-sub', true);

INSERT INTO
	hydra_oauth2_authentication_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, login_session_id, requested_at_audience)
VALUES
	('12-challenge', '12-verifier', '12-client', '12-subject', '12-redirect', false, '12-scope', '12-csrf', NOW(), NOW(), '{}', '12-login-session-id', '12-aud');

INSERT INTO
	hydra_oauth2_consent_request (challenge, verifier, client_id, subject, request_url, skip, requested_scope, csrf, authenticated_at, requested_at, oidc_context, forced_subject_identifier, login_session_id, login_challenge, requested_at_audience, acr, context)
VALUES
	('12-challenge', '12-verifier', '12-client', '12-subject', '12-redirect', false, '12-scope', '12-csrf', NOW(), NOW(), '{}', '12-forced-sub', '12-login-session-id', '12-challenge', '12-aud', '12-acr', '{"foo":"bar"}');

INSERT INTO
	hydra_oauth2_consent_request_handled (challenge, granted_scope, remember, remember_for, error, requested_at, session_access_token, session_id_token, authenticated_at, was_used, granted_at_audience)
VALUES
	('12-challenge', '12-scope', true, 3600, '{}', NOW(), '{}', '{}', NOW(), false, '12-aud');

INSERT INTO
	hydra_oauth2_authentication_request_handled (challenge, subject, remember, remember_for, error, acr, requested_at, authenticated_at, was_used, forced_subject_identifier, context)
VALUES
	('12-challenge', '12-sub', true, 3600, '{}', '1', NOW(), NOW(), false, '12-forced-sub', '{"foo":"bar"}');

INSERT INTO
	hydra_oauth2_obfuscated_authentication_session (subject, client_id, subject_obfuscated)
VALUES
	('12-sub', '12-client', '12-obfuscated');

INSERT INTO
	hydra_oauth2_logout_request (challenge, verifier, subject, sid, client_id, request_url, redir_url, was_used, accepted, rejected, rp_initiated)
VALUES
	('12-challenge', '12-verifier', '12-subject', '12-session-id', '12-client', 'https://request-url/', 'https://redirect-url', false, false, false, false);

INSERT INTO
	hydra_oauth2_logout_request (challenge, verifier, subject, sid, client_id, request_url, redir_url, was_used, accepted, rejected, rp_initiated)
VALUES
	('12-1-challenge', '12-1-verifier', '12-1-subject', '12-1-session-id', NULL, 'https://request-url/', 'https://redirect-url', false, false, false, false);

-- +migrate Down
