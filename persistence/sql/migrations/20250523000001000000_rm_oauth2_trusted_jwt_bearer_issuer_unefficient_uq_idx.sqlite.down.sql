CREATE UNIQUE INDEX IF NOT EXISTS hydra_oauth2_trusted_jwt_bearer_issuer_issuer_subject_key_id_key ON hydra_oauth2_trusted_jwt_bearer_issuer (issuer ASC, subject ASC, key_id ASC, nid ASC);
