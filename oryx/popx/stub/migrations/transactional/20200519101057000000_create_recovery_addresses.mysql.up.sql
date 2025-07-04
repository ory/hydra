CREATE TABLE `identity_recovery_addresses` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`via` VARCHAR (16) NOT NULL,
`value` VARCHAR (400) NOT NULL,
`identity_id` char(36) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`identity_id`) REFERENCES `identities` (`id`) ON DELETE cascade
) ENGINE=InnoDB