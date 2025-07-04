CREATE TABLE `identities` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`traits_schema_id` VARCHAR (2048) NOT NULL,
`traits` JSON NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL
) ENGINE=InnoDB;
CREATE TABLE `identity_credential_types` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`name` VARCHAR (32) NOT NULL
) ENGINE=InnoDB;
CREATE UNIQUE INDEX `identity_credential_types_name_idx` ON `identity_credential_types` (`name`);
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
) ENGINE=InnoDB;
CREATE TABLE `identity_credential_identifiers` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`identifier` VARCHAR (255) NOT NULL,
`identity_credential_id` char(36) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`identity_credential_id`) REFERENCES `identity_credentials` (`id`) ON DELETE cascade
) ENGINE=InnoDB;
CREATE UNIQUE INDEX `identity_credential_identifiers_identifier_idx` ON `identity_credential_identifiers` (`identifier`);