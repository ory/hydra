-- +migrate Up
ALTER TABLE hydra_jwk DROP PRIMARY KEY;
CREATE UNIQUE INDEX hydra_jwk_idx_id_uq ON hydra_jwk (sid, kid);
ALTER TABLE hydra_jwk ADD pk INT UNSIGNED AUTO_INCREMENT PRIMARY KEY;

-- +migrate Down
ALTER TABLE hydra_jwk DROP COLUMN pk;
ALTER TABLE hydra_jwk DROP INDEX hydra_jwk_idx_id_uq;
ALTER TABLE hydra_jwk ADD PRIMARY KEY (sid, kid);
