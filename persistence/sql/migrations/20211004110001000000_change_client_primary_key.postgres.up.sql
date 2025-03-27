ALTER TABLE hydra_client RENAME pk TO pk_deprecated;
ALTER TABLE hydra_client ADD COLUMN pk UUID DEFAULT gen_random_uuid();
ALTER TABLE hydra_client ALTER pk DROP DEFAULT;
ALTER TABLE hydra_client DROP CONSTRAINT hydra_client_pkey;
ALTER TABLE hydra_client ADD PRIMARY KEY (pk);
