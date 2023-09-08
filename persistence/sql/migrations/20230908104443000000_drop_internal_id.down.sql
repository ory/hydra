ALTER TABLE hydra_client DROP CONSTRAINT hydra_client_pkey;

ALTER TABLE hydra_client ADD COLUMN pk UUID NOT NULL DEFAULT gen_random_uuid();

ALTER TABLE hydra_client ADD PRIMARY KEY (pk);
