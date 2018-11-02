-- +migrate Up
ALTER TABLE hydra_client DROP PRIMARY KEY;
CREATE UNIQUE INDEX hydra_client_idx_id_uq ON hydra_client (id);
ALTER TABLE hydra_client ADD pk INT UNSIGNED AUTO_INCREMENT PRIMARY KEY;

-- +migrate Down
ALTER TABLE hydra_client DROP COLUMN pk;
ALTER TABLE hydra_client DROP INDEX hydra_client_idx_id_uq;
ALTER TABLE hydra_client ADD PRIMARY KEY (id);
