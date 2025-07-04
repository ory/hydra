ALTER TABLE `selfservice_login_flow_methods` CHANGE `selfservice_login_request_id` `selfservice_login_flow_id` char(36) NOT NULL;
ALTER TABLE `selfservice_registration_flow_methods` CHANGE `selfservice_registration_request_id` `selfservice_registration_flow_id` char(36) NOT NULL;
ALTER TABLE `selfservice_recovery_flow_methods` CHANGE `selfservice_recovery_request_id` `selfservice_recovery_flow_id` char(36) NOT NULL;
ALTER TABLE `selfservice_settings_flow_methods` CHANGE `selfservice_settings_request_id` `selfservice_settings_flow_id` char(36) NOT NULL;