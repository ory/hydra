ALTER TABLE hydra_jwk DROP CONSTRAINT "primary";
ALTER TABLE hydra_jwk ADD CONSTRAINT "hydra_jwk_pkey" PRIMARY KEY (pk);
