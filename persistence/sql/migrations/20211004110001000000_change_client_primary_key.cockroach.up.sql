ALTER TABLE hydra_client RENAME pk TO pk_deprecated;
ALTER TABLE hydra_client ADD pk UUID NOT NULL DEFAULT gen_random_uuid();
