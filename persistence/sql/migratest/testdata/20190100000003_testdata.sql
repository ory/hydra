INSERT INTO hydra_client
(id, client_name, client_secret, redirect_uris, grant_types, response_types, scope, owner, policy_uri, tos_uri, client_uri, logo_uri, contacts, public, client_secret_expires_at, sector_identifier_uri, jwks, jwks_uri, request_uris, token_endpoint_auth_method, request_object_signing_alg, userinfo_signed_response_alg)
VALUES
('client-0003', 'Client 0003', 'secret-0003', 'http://redirect/0003_1', 'grant-0003_1', 'response-0003_1', 'scope-0003', 'owner-0003', 'http://policy/0003', 'http://tos/0003', 'http://client/0003', 'http://logo/0003', 'contact-0003_1', true, 0, NULL, NULL, NULL, 'http://request/0003_1', 'token_auth-0003', 'r_alg-0003', 'u_alg-0003');
