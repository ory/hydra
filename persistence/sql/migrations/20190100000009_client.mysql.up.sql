ALTER TABLE hydra_client DROP PRIMARY KEY;
CREATE UNIQUE INDEX hydra_client_idx_id_uq ON hydra_client (id);
ALTER TABLE hydra_client ADD pk INT UNSIGNED AUTO_INCREMENT PRIMARY KEY;
