ALTER TABLE hydra_jwk
    RENAME pk TO pk_deprecated,
    ADD COLUMN pk UUID NOT NULL DEFAULT gen_random_uuid();
