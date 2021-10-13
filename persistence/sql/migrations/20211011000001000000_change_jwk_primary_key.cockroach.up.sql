ALTER TABLE hydra_jwk RENAME pk TO pk_deprecated;
ALTER TABLE hydra_jwk ADD pk UUID NOT NULL DEFAULT gen_random_uuid();
