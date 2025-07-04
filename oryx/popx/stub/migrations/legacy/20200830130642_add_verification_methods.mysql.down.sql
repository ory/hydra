ALTER TABLE `selfservice_verification_flows` ADD COLUMN `form` JSON;
UPDATE selfservice_verification_flows SET form=(SELECT * FROM (SELECT m.config FROM selfservice_verification_flows AS r INNER JOIN selfservice_verification_flow_methods AS m ON r.id=m.selfservice_verification_flow_id) as t);
ALTER TABLE `selfservice_verification_flows` MODIFY `form` JSON;
DROP TABLE `selfservice_verification_flow_methods`;
ALTER TABLE `selfservice_verification_flows` DROP COLUMN `active_method`;
ALTER TABLE `selfservice_verification_flows` DROP COLUMN `state`;
ALTER TABLE `selfservice_verification_flows` ADD COLUMN `via` VARCHAR (16) NOT NULL DEFAULT 'email';
ALTER TABLE `selfservice_verification_flows` ADD COLUMN `success` bool NOT NULL DEFAULT FALSE;