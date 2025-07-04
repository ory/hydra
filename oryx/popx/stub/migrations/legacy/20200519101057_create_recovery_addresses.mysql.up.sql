CREATE TABLE `identity_recovery_addresses` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`via` VARCHAR (16) NOT NULL,
`value` VARCHAR (400) NOT NULL,
`identity_id` char(36) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`identity_id`) REFERENCES `identities` (`id`) ON DELETE cascade
) ENGINE=InnoDB;
CREATE UNIQUE INDEX `identity_recovery_addresses_status_via_uq_idx` ON `identity_recovery_addresses` (`via`, `value`);
CREATE INDEX `identity_recovery_addresses_status_via_idx` ON `identity_recovery_addresses` (`via`, `value`);
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
) ENGINE=InnoDB;
CREATE TABLE `selfservice_recovery_request_methods` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`method` VARCHAR (32) NOT NULL,
`config` JSON NOT NULL,
`selfservice_recovery_request_id` char(36) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`selfservice_recovery_request_id`) REFERENCES `selfservice_recovery_requests` (`id`) ON DELETE cascade
) ENGINE=InnoDB;
CREATE TABLE `identity_recovery_tokens` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`token` VARCHAR (64) NOT NULL,
`used` bool NOT NULL DEFAULT false,
`used_at` DATETIME,
`identity_recovery_address_id` char(36) NOT NULL,
`selfservice_recovery_request_id` char(36) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`identity_recovery_address_id`) REFERENCES `identity_recovery_addresses` (`id`) ON DELETE cascade,
FOREIGN KEY (`selfservice_recovery_request_id`) REFERENCES `selfservice_recovery_requests` (`id`) ON DELETE cascade
) ENGINE=InnoDB;
CREATE UNIQUE INDEX `identity_recovery_addresses_code_uq_idx` ON `identity_recovery_tokens` (`token`);
CREATE INDEX `identity_recovery_addresses_code_idx` ON `identity_recovery_tokens` (`token`);