ALTER TABLE hydra_client CHANGE COLUMN pk pk_deprecated INT UNSIGNED;
ALTER TABLE hydra_client ADD COLUMN pk_new CHAR(36);
UPDATE hydra_client SET pk_new = (SELECT uuid());
ALTER TABLE hydra_client ALTER pk_new DROP DEFAULT;
ALTER TABLE hydra_client DROP PRIMARY KEY, ADD PRIMARY KEY (pk_new);
ALTER TABLE hydra_client ADD KEY (pk_deprecated);
ALTER TABLE hydra_client CHANGE COLUMN pk_deprecated pk_deprecated INT UNSIGNED AUTO_INCREMENT;
