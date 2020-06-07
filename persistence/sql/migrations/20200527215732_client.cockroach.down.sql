UPDATE hydra_client SET token_endpoint_auth_signing_alg = '' WHERE token_endpoint_auth_method = 'private_key_jwt' AND token_endpoint_auth_signing_alg = 'RS256';
