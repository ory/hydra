ALTER TABLE hydra_client RENAME pk TO pk_tmp;
ALTER TABLE hydra_client CHANGE COLUMN pk_deprecated pk INT UNSIGNED AUTO_INCREMENT;
ALTER TABLE hydra_client DROP PRIMARY KEY, ADD PRIMARY KEY (pk);
ALTER TABLE hydra_client DROP pk_tmp;
