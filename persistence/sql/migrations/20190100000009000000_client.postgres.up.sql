ALTER TABLE hydra_client DROP CONSTRAINT hydra_client_pkey;
ALTER TABLE hydra_client ADD pk SERIAL;
ALTER TABLE hydra_client ADD PRIMARY KEY (pk);
CREATE UNIQUE INDEX hydra_client_idx_id_uq ON hydra_client (id);
