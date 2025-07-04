CREATE TABLE `identity_credentials` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`config` JSON NOT NULL,
`identity_credential_type_id` char(36) NOT NULL,
`identity_id` char(36) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`identity_id`) REFERENCES `identities` (`id`) ON DELETE cascade,
FOREIGN KEY (`identity_credential_type_id`) REFERENCES `identity_credential_types` (`id`) ON DELETE cascade
) ENGINE=InnoDB