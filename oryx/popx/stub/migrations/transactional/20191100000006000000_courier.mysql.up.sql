CREATE TABLE `courier_messages` (
`id` char(36) NOT NULL,
PRIMARY KEY(`id`),
`type` INTEGER NOT NULL,
`status` INTEGER NOT NULL,
`body` VARCHAR (255) NOT NULL,
`subject` VARCHAR (255) NOT NULL,
`recipient` VARCHAR (255) NOT NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL
) ENGINE=InnoDB;