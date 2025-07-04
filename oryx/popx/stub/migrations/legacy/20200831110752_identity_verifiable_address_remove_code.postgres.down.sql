ALTER TABLE "identity_verifiable_addresses" ADD COLUMN "code" VARCHAR (32);
ALTER TABLE "identity_verifiable_addresses" ADD COLUMN "expires_at" timestamp;
UPDATE identity_verifiable_addresses SET code = substr(md5(random()::text), 0, 32) WHERE code IS NULL;
UPDATE identity_verifiable_addresses SET expires_at = CURRENT_TIMESTAMP WHERE expires_at IS NULL;
ALTER TABLE "identity_verifiable_addresses" ALTER COLUMN "code" TYPE VARCHAR (32), ALTER COLUMN "code" SET NOT NULL;
ALTER TABLE "identity_verifiable_addresses" ALTER COLUMN "expires_at" TYPE timestamp, ALTER COLUMN "expires_at" DROP NOT NULL;
CREATE UNIQUE INDEX "identity_verifiable_addresses_code_uq_idx" ON "identity_verifiable_addresses" (code);
CREATE INDEX "identity_verifiable_addresses_code_idx" ON "identity_verifiable_addresses" (code);