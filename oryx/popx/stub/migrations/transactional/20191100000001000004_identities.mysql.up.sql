CREATE TABLE `identity_credential_identifiers` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`identifier` VARCHAR (255) NOT NULL,
`identity_credential_id` char(36) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`identity_credential_id`) REFERENCES `identity_credentials` (`id`) ON DELETE cascade
) ENGINE=InnoDB