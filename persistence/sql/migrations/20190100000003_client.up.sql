ALTER TABLE hydra_client ADD sector_identifier_uri TEXT;
ALTER TABLE hydra_client ADD jwks TEXT;
ALTER TABLE hydra_client ADD jwks_uri TEXT;
ALTER TABLE hydra_client ADD request_uris TEXT;
ALTER TABLE hydra_client ADD token_endpoint_auth_method VARCHAR(25) NOT NULL DEFAULT '';
ALTER TABLE hydra_client ADD request_object_signing_alg  VARCHAR(10) NOT NULL DEFAULT '';
ALTER TABLE hydra_client ADD userinfo_signed_response_alg VARCHAR(10) NOT NULL DEFAULT '';
