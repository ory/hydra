ALTER TABLE hydra_jwk RENAME pk TO pk_deprecated;
ALTER TABLE hydra_jwk ADD COLUMN pk UUID DEFAULT gen_random_uuid();
ALTER TABLE hydra_jwk ALTER pk DROP DEFAULT;
ALTER TABLE hydra_jwk DROP CONSTRAINT hydra_jwk_pkey;
ALTER TABLE hydra_jwk ADD PRIMARY KEY (pk);
