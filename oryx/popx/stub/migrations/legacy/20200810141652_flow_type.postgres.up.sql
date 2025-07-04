ALTER TABLE "selfservice_login_requests" ADD COLUMN "type" VARCHAR (16) NOT NULL DEFAULT 'browser';
ALTER TABLE "selfservice_registration_requests" ADD COLUMN "type" VARCHAR (16) NOT NULL DEFAULT 'browser';
ALTER TABLE "selfservice_settings_requests" ADD COLUMN "type" VARCHAR (16) NOT NULL DEFAULT 'browser';
ALTER TABLE "selfservice_recovery_requests" ADD COLUMN "type" VARCHAR (16) NOT NULL DEFAULT 'browser';
ALTER TABLE "selfservice_verification_requests" ADD COLUMN "type" VARCHAR (16) NOT NULL DEFAULT 'browser';