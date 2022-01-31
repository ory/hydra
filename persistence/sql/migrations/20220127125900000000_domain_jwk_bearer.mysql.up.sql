
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    ADD COLUMN allowed_domain VARCHAR(255) NOT NULL DEFAULT '',
    DROP INDEX issuer;

ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    ADD CONSTRAINT hydra_oauth2_trusted_jwt_bearer_issuer_issuer_subject_key_idx UNIQUE (issuer(128), subject(128), allowed_domain(128), key_id);
