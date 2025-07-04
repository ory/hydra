ALTER TABLE `selfservice_profile_management_requests` ADD COLUMN `form` JSON;
UPDATE selfservice_profile_management_requests SET form=(SELECT * FROM (SELECT m.config FROM selfservice_profile_management_requests AS r INNER JOIN selfservice_profile_management_request_methods AS m ON r.id=m.selfservice_profile_management_request_id) as t);
ALTER TABLE `selfservice_profile_management_requests` MODIFY `form` JSON;
DROP TABLE `selfservice_profile_management_request_methods`;
ALTER TABLE `selfservice_profile_management_requests` DROP COLUMN `active_method`;