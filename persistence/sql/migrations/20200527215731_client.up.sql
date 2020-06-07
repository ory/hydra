ALTER TABLE hydra_client ADD token_endpoint_auth_signing_alg VARCHAR(10) NOT NULL DEFAULT '';
UPDATE hydra_client SET token_endpoint_auth_signing_alg = 'RS256' WHERE token_endpoint_auth_method = 'private_key_jwt';
