CREATE TABLE "selfservice_profile_management_request_methods" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"method" VARCHAR (32) NOT NULL,
"selfservice_profile_management_request_id" UUID NOT NULL,
"config" json NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL
);COMMIT TRANSACTION;BEGIN TRANSACTION;
ALTER TABLE "selfservice_profile_management_requests" ADD COLUMN "active_method" VARCHAR (32);COMMIT TRANSACTION;BEGIN TRANSACTION;
INSERT INTO selfservice_profile_management_request_methods (id, method, selfservice_profile_management_request_id, config) SELECT id, 'traits', id, form FROM selfservice_profile_management_requests;
ALTER TABLE "selfservice_profile_management_requests" DROP COLUMN "form";COMMIT TRANSACTION;BEGIN TRANSACTION;