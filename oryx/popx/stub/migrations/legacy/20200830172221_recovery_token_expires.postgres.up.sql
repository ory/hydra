ALTER TABLE "identity_recovery_tokens" ADD COLUMN "expires_at" timestamp NOT NULL DEFAULT '2000-01-01 00:00:00';
ALTER TABLE "identity_recovery_tokens" ADD COLUMN "issued_at" timestamp NOT NULL DEFAULT '2000-01-01 00:00:00';
ALTER TABLE "identity_recovery_tokens" ALTER COLUMN "selfservice_recovery_flow_id" TYPE UUID, ALTER COLUMN "selfservice_recovery_flow_id" DROP NOT NULL;