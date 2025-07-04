ALTER TABLE "selfservice_profile_management_requests" ADD COLUMN "form" json NOT NULL DEFAULT '{}';COMMIT TRANSACTION;BEGIN TRANSACTION;
DROP TABLE "selfservice_profile_management_request_methods";COMMIT TRANSACTION;BEGIN TRANSACTION;
ALTER TABLE "selfservice_profile_management_requests" DROP COLUMN "active_method";COMMIT TRANSACTION;BEGIN TRANSACTION;