CREATE TABLE `selfservice_login_request_methods` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`method` VARCHAR (32) NOT NULL,
`selfservice_login_request_id` char(36) NOT NULL,
`config` JSON NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
FOREIGN KEY (`selfservice_login_request_id`) REFERENCES `selfservice_login_requests` (`id`) ON DELETE cascade
) ENGINE=InnoDB