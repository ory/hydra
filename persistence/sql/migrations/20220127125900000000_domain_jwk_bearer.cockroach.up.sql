DROP INDEX hydra_oauth2_trusted_jwt_bearer_issuer_issuer_subject_key_id_key CASCADE;

ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    ADD COLUMN domain VARCHAR(255) NOT NULL DEFAULT '';

CREATE UNIQUE INDEX ON hydra_oauth2_trusted_jwt_bearer_issuer(issuer, subject, domain, key_id);