-- +migrate Up
ALTER TABLE hydra_client ADD sector_identifier_uri TEXT;
ALTER TABLE hydra_client ADD jwks TEXT;
ALTER TABLE hydra_client ADD jwks_uri TEXT;
ALTER TABLE hydra_client ADD request_uris TEXT;
ALTER TABLE hydra_client ADD token_endpoint_auth_method VARCHAR(25) NOT NULL DEFAULT '';
ALTER TABLE hydra_client ADD request_object_signing_alg  VARCHAR(10) NOT NULL DEFAULT '';
ALTER TABLE hydra_client ADD userinfo_signed_response_alg VARCHAR(10) NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE hydra_client DROP COLUMN sector_identifier_uri;
ALTER TABLE hydra_client DROP COLUMN jwks;
ALTER TABLE hydra_client DROP COLUMN jwks_uri;
ALTER TABLE hydra_client DROP COLUMN token_endpoint_auth_method;
ALTER TABLE hydra_client DROP COLUMN request_uris;
ALTER TABLE hydra_client DROP COLUMN request_object_signing_alg;
ALTER TABLE hydra_client DROP COLUMN userinfo_signed_response_alg;
