ALTER TABLE hydra_jwk RENAME pk TO pk_tmp;
ALTER TABLE hydra_jwk RENAME pk_deprecated TO pk;
ALTER TABLE hydra_jwk DROP CONSTRAINT hydra_jwk_pkey;
ALTER TABLE hydra_jwk ADD PRIMARY KEY (pk);
ALTER TABLE hydra_jwk DROP pk_tmp;
