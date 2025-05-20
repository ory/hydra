ALTER TABLE hydra_jwk
    DROP COLUMN pk,
    RENAME pk_deprecated TO pk;
