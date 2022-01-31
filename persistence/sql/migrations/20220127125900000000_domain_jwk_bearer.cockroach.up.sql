DROP INDEX hydra_oauth2_trusted_jwt_bearer_issuer_issuer_subject_key_id_key CASCADE;

ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    ADD COLUMN allowed_domain VARCHAR(255) NOT NULL DEFAULT '';

CREATE UNIQUE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_issuer_subject_allowed_domain_key_id_key
    ON hydra_oauth2_trusted_jwt_bearer_issuer(issuer, subject, allowed_domain, key_id);
