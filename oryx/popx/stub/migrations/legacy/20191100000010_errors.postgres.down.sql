UPDATE selfservice_errors SET seen_at = '1980-01-01 00:00:00' WHERE seen_at = NULL;
ALTER TABLE "selfservice_errors" ALTER COLUMN "seen_at" TYPE timestamp, ALTER COLUMN "seen_at" DROP NOT NULL;