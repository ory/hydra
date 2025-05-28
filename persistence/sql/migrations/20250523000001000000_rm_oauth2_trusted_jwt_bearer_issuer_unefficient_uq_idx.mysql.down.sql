CREATE UNIQUE INDEX issuer ON hydra_oauth2_trusted_jwt_bearer_issuer (issuer, subject, key_id, nid);
