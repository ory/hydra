ALTER TABLE hydra_jwk
    DROP CONSTRAINT "primary",
    ADD CONSTRAINT "primary" PRIMARY KEY (pk_deprecated)
