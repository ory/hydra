ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    ADD COLUMN allow_any_subject BOOL NOT NULL DEFAULT FALSE;
