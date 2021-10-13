ALTER TABLE hydra_client DROP PRIMARY KEY, ADD PRIMARY KEY (pk_deprecated);
ALTER TABLE hydra_client DROP KEY pk_deprecated;
ALTER TABLE hydra_client DROP pk;
ALTER TABLE hydra_client CHANGE COLUMN pk_deprecated pk INT UNSIGNED AUTO_INCREMENT;
