DELETE FROM identity_recovery_tokens WHERE selfservice_recovery_flow_id IS NULL;
ALTER TABLE "identity_recovery_tokens" ALTER COLUMN "selfservice_recovery_flow_id" TYPE UUID, ALTER COLUMN "selfservice_recovery_flow_id" SET NOT NULL;
ALTER TABLE "identity_recovery_tokens" DROP COLUMN "expires_at";
ALTER TABLE "identity_recovery_tokens" DROP COLUMN "issued_at";