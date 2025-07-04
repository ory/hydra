CREATE TABLE `sessions` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`issued_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`expires_at` DATETIME NOT NULL,
`authenticated_at` DATETIME NOT NULL,
`identity_id` char(36) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`identity_id`) REFERENCES `identities` (`id`) ON DELETE cascade
) ENGINE=InnoDB;