ALTER TABLE hydra_jwk DROP CONSTRAINT hydra_jwk_pkey;
ALTER TABLE hydra_jwk ADD pk SERIAL;
ALTER TABLE hydra_jwk ADD PRIMARY KEY (pk);
CREATE UNIQUE INDEX hydra_jwk_idx_id_uq ON hydra_jwk (sid, kid);
