DROP INDEX hydra_oauth2_trusted_jwt_bearer_issuer_issuer_subject_domain_key_id_key CASCADE;

ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    DROP COLUMN domain;

CREATE UNIQUE INDEX ON hydra_oauth2_trusted_jwt_bearer_issuer(issuer, subject, key_id);
