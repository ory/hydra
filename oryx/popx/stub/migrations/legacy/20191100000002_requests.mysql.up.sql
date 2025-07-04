CREATE TABLE `selfservice_login_requests` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`request_url` VARCHAR (2048) NOT NULL,
`issued_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`expires_at` DATETIME NOT NULL,
`active_method` VARCHAR (32) NOT NULL,
`csrf_token` VARCHAR (255) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL
) ENGINE=InnoDB;
CREATE TABLE `selfservice_login_request_methods` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`method` VARCHAR (32) NOT NULL,
`selfservice_login_request_id` char(36) NOT NULL,
`config` JSON NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`selfservice_login_request_id`) REFERENCES `selfservice_login_requests` (`id`) ON DELETE cascade
) ENGINE=InnoDB;
CREATE TABLE `selfservice_registration_requests` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`request_url` VARCHAR (2048) NOT NULL,
`issued_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`expires_at` DATETIME NOT NULL,
`active_method` VARCHAR (32) NOT NULL,
`csrf_token` VARCHAR (255) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL
) ENGINE=InnoDB;
CREATE TABLE `selfservice_registration_request_methods` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`method` VARCHAR (32) NOT NULL,
`selfservice_registration_request_id` char(36) NOT NULL,
`config` JSON NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`selfservice_registration_request_id`) REFERENCES `selfservice_registration_requests` (`id`) ON DELETE cascade
) ENGINE=InnoDB;
CREATE TABLE `selfservice_profile_management_requests` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`request_url` VARCHAR (2048) NOT NULL,
`issued_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`expires_at` DATETIME NOT NULL,
`form` JSON NOT NULL,
`update_successful` bool NOT NULL,
`identity_id` char(36) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`identity_id`) REFERENCES `identities` (`id`) ON DELETE cascade
) ENGINE=InnoDB;