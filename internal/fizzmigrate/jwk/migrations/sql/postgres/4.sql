-- +migrate Up
ALTER TABLE hydra_jwk DROP CONSTRAINT hydra_jwk_pkey;
ALTER TABLE hydra_jwk ADD pk SERIAL;
ALTER TABLE hydra_jwk ADD PRIMARY KEY (pk);
CREATE UNIQUE INDEX hydra_jwk_idx_id_uq ON hydra_jwk (sid, kid);

-- +migrate Down
ALTER TABLE hydra_jwk DROP CONSTRAINT hydra_jwk_pkey;
ALTER TABLE hydra_jwk DROP COLUMN pk;
DROP INDEX hydra_jwk_idx_id_uq;
ALTER TABLE hydra_jwk ADD PRIMARY KEY (sid, kid);
