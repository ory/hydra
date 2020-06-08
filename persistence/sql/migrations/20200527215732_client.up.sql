-- complement for migration 20200527215731_client.up
-- due to cockroach being unable to handle a scheme modification and update statement in a single migration transaction
UPDATE hydra_client SET token_endpoint_auth_signing_alg = 'RS256' WHERE token_endpoint_auth_method = 'private_key_jwt';
