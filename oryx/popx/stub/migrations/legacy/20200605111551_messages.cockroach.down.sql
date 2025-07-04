ALTER TABLE "selfservice_verification_requests" DROP COLUMN "messages";COMMIT TRANSACTION;BEGIN TRANSACTION;
ALTER TABLE "selfservice_login_requests" DROP COLUMN "messages";COMMIT TRANSACTION;BEGIN TRANSACTION;
ALTER TABLE "selfservice_registration_requests" DROP COLUMN "messages";COMMIT TRANSACTION;BEGIN TRANSACTION;