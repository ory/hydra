ALTER TABLE `selfservice_settings_request_methods` CHANGE `selfservice_settings_request_id` `selfservice_profile_management_request_id` char(36) NOT NULL;
ALTER TABLE `selfservice_settings_request_methods` RENAME TO `selfservice_profile_management_request_methods`;
ALTER TABLE `selfservice_settings_requests` RENAME TO `selfservice_profile_management_requests`;