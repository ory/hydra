CREATE TABLE `networks` (
  `id` char(36) NOT NULL,
  PRIMARY KEY(`id`),
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL
);

INSERT INTO networks (id, created_at, updated_at) VALUES ((SELECT LOWER(CONCAT(
    HEX(RANDOM_BYTES(4)),
    '-', HEX(RANDOM_BYTES(2)),
    '-4', SUBSTR(HEX(RANDOM_BYTES(2)), 2, 3),
    '-', CONCAT(HEX(FLOOR(ASCII(RANDOM_BYTES(1)) / 64)+8),SUBSTR(HEX(RANDOM_BYTES(2)), 2, 3)),
    '-', HEX(RANDOM_BYTES(6))
))), '2013-10-07 08:23:19', '2013-10-07 08:23:19');
