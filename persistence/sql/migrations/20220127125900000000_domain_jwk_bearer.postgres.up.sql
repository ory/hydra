
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    ADD COLUMN allowed_domain VARCHAR(255) NOT NULL DEFAULT '',
    DROP CONSTRAINT hydra_oauth2_trusted_jwt_bearer_issue_issuer_subject_key_id_key;

ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    ADD CONSTRAINT hydra_oauth2_trusted_jwt_bearer_issuer_issuer_subject_key_idx UNIQUE (issuer, subject, allowed_domain, key_id);
