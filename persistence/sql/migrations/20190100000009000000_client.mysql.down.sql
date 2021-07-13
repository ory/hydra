ALTER TABLE hydra_client DROP COLUMN pk;
ALTER TABLE hydra_client DROP INDEX hydra_client_idx_id_uq;
ALTER TABLE hydra_client ADD PRIMARY KEY (id);
