CREATE TABLE `networks` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL
) ENGINE=InnoDB;