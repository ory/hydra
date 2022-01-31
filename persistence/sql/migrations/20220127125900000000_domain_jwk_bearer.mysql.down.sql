
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    DROP INDEX hydra_oauth2_trusted_jwt_bearer_issuer_issuer_subject_key_idx,
    ADD CONSTRAINT issuer UNIQUE (issuer, subject, key_id);

ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    DROP COLUMN allowed_domain;
