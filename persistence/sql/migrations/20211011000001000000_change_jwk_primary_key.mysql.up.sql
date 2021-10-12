ALTER TABLE hydra_jwk CHANGE COLUMN pk pk_deprecated INT UNSIGNED;
ALTER TABLE hydra_jwk ADD COLUMN pk CHAR(36);
-- UUIDv4 generation based on https://stackoverflow.com/a/66868340/12723442
UPDATE hydra_jwk SET pk = (SELECT LOWER(CONCAT(
    HEX(RANDOM_BYTES(4)),
    '-', HEX(RANDOM_BYTES(2)),
    '-4', SUBSTR(HEX(RANDOM_BYTES(2)), 2, 3),
    '-', CONCAT(HEX(FLOOR(ASCII(RANDOM_BYTES(1)) / 64)+8),SUBSTR(HEX(RANDOM_BYTES(2)), 2, 3)),
    '-', HEX(RANDOM_BYTES(6))
)));
ALTER TABLE hydra_jwk ALTER pk DROP DEFAULT;
ALTER TABLE hydra_jwk DROP PRIMARY KEY, ADD PRIMARY KEY (pk);
ALTER TABLE hydra_jwk ADD KEY (pk_deprecated);
ALTER TABLE hydra_jwk CHANGE COLUMN pk_deprecated pk_deprecated INT UNSIGNED AUTO_INCREMENT;
