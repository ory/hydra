ALTER TABLE hydra_client
    DROP CONSTRAINT "primary",
    ADD CONSTRAINT "hydra_client_pkey" PRIMARY KEY (pk_deprecated);
