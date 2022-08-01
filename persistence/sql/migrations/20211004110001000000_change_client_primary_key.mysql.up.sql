ALTER TABLE hydra_client CHANGE COLUMN pk pk_deprecated INT UNSIGNED;
ALTER TABLE hydra_client ADD COLUMN pk CHAR(36);
-- UUIDv4 generation based on https://stackoverflow.com/a/66868340/12723442
UPDATE hydra_client SET pk = (SELECT LOWER(CONCAT(
    HEX(RANDOM_BYTES(4)),
    '-', HEX(RANDOM_BYTES(2)),
    '-4', SUBSTR(HEX(RANDOM_BYTES(2)), 2, 3),
    '-', CONCAT(HEX(FLOOR(ASCII(RANDOM_BYTES(1)) / 64)+8),SUBSTR(HEX(RANDOM_BYTES(2)), 2, 3)),
    '-', HEX(RANDOM_BYTES(6))
)));
ALTER TABLE hydra_client DROP PRIMARY KEY, ADD PRIMARY KEY (pk);
ALTER TABLE hydra_client ADD KEY (pk_deprecated);
ALTER TABLE hydra_client CHANGE COLUMN pk_deprecated pk_deprecated INT UNSIGNED;
