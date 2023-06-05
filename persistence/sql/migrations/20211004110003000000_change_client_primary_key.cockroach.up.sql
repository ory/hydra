ALTER TABLE hydra_client DROP CONSTRAINT IF EXISTS "primary";
ALTER TABLE hydra_client ADD CONSTRAINT "hydra_client_pkey" IF NOT EXISTS PRIMARY KEY (pk);
