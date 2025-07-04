CREATE TABLE `continuity_containers` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`identity_id` char(36),
`name` VARCHAR (255) NOT NULL,
`payload` JSON,
`expires_at` DATETIME NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`identity_id`) REFERENCES `identities` (`id`) ON DELETE cascade
) ENGINE=InnoDB;