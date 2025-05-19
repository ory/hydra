ALTER TABLE hydra_client 
    RENAME pk TO pk_deprecated,
    ADD COLUMN pk UUID NOT NULL DEFAULT gen_random_uuid();
