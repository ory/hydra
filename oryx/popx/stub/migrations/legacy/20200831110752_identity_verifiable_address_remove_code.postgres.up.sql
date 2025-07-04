DROP INDEX "identity_verifiable_addresses_code_uq_idx";
DROP INDEX "identity_verifiable_addresses_code_idx";
ALTER TABLE "identity_verifiable_addresses" DROP COLUMN "code";
ALTER TABLE "identity_verifiable_addresses" DROP COLUMN "expires_at";