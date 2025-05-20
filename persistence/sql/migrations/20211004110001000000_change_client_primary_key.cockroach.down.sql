ALTER TABLE hydra_client "hydra_client_pkey"
    DROP COLUMN pk,
    RENAME COLUMN pk_deprecated TO pk;
