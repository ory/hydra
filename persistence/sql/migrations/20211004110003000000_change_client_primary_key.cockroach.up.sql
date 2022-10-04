ALTER TABLE hydra_client DROP CONSTRAINT "primary";
ALTER TABLE hydra_client ADD CONSTRAINT "hydra_client_pkey" PRIMARY KEY (pk);
