ALTER TABLE "selfservice_settings_request_methods" RENAME COLUMN "selfservice_settings_request_id" TO "selfservice_profile_management_request_id";
ALTER TABLE "selfservice_settings_request_methods" RENAME TO "selfservice_profile_management_request_methods";
ALTER TABLE "selfservice_settings_requests" RENAME TO "selfservice_profile_management_requests";