ALTER TABLE hydra_client CHANGE COLUMN pk pk_deprecated INT UNSIGNED;
ALTER TABLE hydra_client ADD COLUMN pk CHAR(36);
UPDATE hydra_client SET pk = (SELECT uuid());
ALTER TABLE hydra_client ALTER pk DROP DEFAULT;
ALTER TABLE hydra_client DROP PRIMARY KEY, ADD PRIMARY KEY (pk);
ALTER TABLE hydra_client ADD KEY (pk_deprecated);
ALTER TABLE hydra_client CHANGE COLUMN pk_deprecated pk_deprecated INT UNSIGNED AUTO_INCREMENT;
