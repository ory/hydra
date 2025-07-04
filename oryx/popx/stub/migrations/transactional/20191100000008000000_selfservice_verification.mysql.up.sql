CREATE TABLE `identity_verifiable_addresses` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`code` VARCHAR (32) NOT NULL,
`status` VARCHAR (16) NOT NULL,
`via` VARCHAR (16) NOT NULL,
`verified` bool NOT NULL,
`value` VARCHAR (400) NOT NULL,
`verified_at` DATETIME,
`expires_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`identity_id` char(36) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`identity_id`) REFERENCES `identities` (`id`) ON DELETE cascade
) ENGINE=InnoDB