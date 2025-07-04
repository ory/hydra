CREATE TABLE `selfservice_errors` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`errors` JSON NOT NULL,
`seen_at` DATETIME NOT NULL,
`was_seen` bool NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL
) ENGINE=InnoDB;