ALTER TABLE hydra_client DROP CONSTRAINT "primary";
ALTER TABLE hydra_client ADD CONSTRAINT "primary" PRIMARY KEY (pk);
