ALTER TABLE `selfservice_settings_requests` DROP COLUMN `state`;
ALTER TABLE `selfservice_settings_requests` ADD COLUMN `update_successful` bool NOT NULL DEFAULT false;