ALTER TABLE hydra_jwk DROP CONSTRAINT "primary";
ALTER TABLE hydra_jwk ADD CONSTRAINT "hydra_client_pkey" PRIMARY KEY (pk);
