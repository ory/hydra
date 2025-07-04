CREATE TABLE `identities` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`traits_schema_id` VARCHAR (2048) NOT NULL,
`traits` JSON NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL
) ENGINE=InnoDB