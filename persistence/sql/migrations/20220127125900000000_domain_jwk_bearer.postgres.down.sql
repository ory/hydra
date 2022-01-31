
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    DROP CONSTRAINT hydra_oauth2_trusted_jwt_bearer_issuer_issuer_subject_key_idx,
    ADD CONSTRAINT hydra_oauth2_trusted_jwt_bearer_issue_issuer_subject_key_id_key UNIQUE (issuer, subject, key_id);

ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer
    DROP COLUMN allowed_domain;
