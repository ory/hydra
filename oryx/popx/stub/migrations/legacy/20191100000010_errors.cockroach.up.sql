ALTER TABLE "selfservice_errors" RENAME COLUMN "seen_at" TO "_seen_at_tmp";COMMIT TRANSACTION;BEGIN TRANSACTION;
ALTER TABLE "selfservice_errors" ADD COLUMN "seen_at" timestamp;COMMIT TRANSACTION;BEGIN TRANSACTION;
UPDATE "selfservice_errors" SET "seen_at" = "_seen_at_tmp";COMMIT TRANSACTION;BEGIN TRANSACTION;
ALTER TABLE "selfservice_errors" DROP COLUMN "_seen_at_tmp";COMMIT TRANSACTION;BEGIN TRANSACTION;