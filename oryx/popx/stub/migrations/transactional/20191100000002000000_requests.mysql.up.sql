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
) ENGINE=InnoDB