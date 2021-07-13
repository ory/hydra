ALTER TABLE hydra_jwk DROP COLUMN pk;
ALTER TABLE hydra_jwk DROP INDEX hydra_jwk_idx_id_uq;
ALTER TABLE hydra_jwk ADD PRIMARY KEY (sid, kid);
