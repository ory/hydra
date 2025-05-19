ALTER TABLE hydra_jwk
    DROP CONSTRAINT "primary",
    ADD CONSTRAINT "hydra_jwk_pkey" PRIMARY KEY (pk);
