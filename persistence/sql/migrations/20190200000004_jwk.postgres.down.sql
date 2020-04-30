ALTER TABLE hydra_jwk DROP CONSTRAINT hydra_jwk_pkey;
ALTER TABLE hydra_jwk DROP COLUMN pk;
DROP INDEX hydra_jwk_idx_id_uq;
ALTER TABLE hydra_jwk ADD PRIMARY KEY (sid, kid);
