ALTER TABLE hydra_jwk DROP CONSTRAINT "primary";
ALTER TABLE hydra_jwk ADD CONSTRAINT "primary" PRIMARY KEY (pk);
