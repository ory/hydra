CREATE TABLE `selfservice_recovery_requests` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`request_url` VARCHAR (2048) NOT NULL,
`issued_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`expires_at` DATETIME NOT NULL,
`messages` JSON,
`active_method` VARCHAR (32),
`csrf_token` VARCHAR (255) NOT NULL,
`state` VARCHAR (32) NOT NULL,
`recovered_identity_id` char(36),
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`recovered_identity_id`) REFERENCES `identities` (`id`) ON DELETE cascade
) ENGINE=InnoDB