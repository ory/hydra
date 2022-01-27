
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    ADD COLUMN domain VARCHAR(255) NOT NULL DEFAULT '',
    DROP INDEX issuer;

ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    ADD CONSTRAINT hydra_oauth2_trusted_jwt_bearer_issuer_issuer_subject_key_idx UNIQUE (issuer, subject, domain, key_id);
