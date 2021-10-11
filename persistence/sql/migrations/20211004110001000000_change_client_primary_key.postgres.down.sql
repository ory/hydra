ALTER TABLE hydra_client RENAME pk TO pk_tmp;
ALTER TABLE hydra_client RENAME pk_deprecated TO pk;
ALTER TABLE hydra_client DROP CONSTRAINT hydra_client_pkey;
ALTER TABLE hydra_client ADD PRIMARY KEY (pk);
ALTER TABLE hydra_client DROP pk_tmp;
