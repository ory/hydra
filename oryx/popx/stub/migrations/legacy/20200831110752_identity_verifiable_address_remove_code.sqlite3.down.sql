ALTER TABLE "identity_verifiable_addresses" ADD COLUMN "code" TEXT;
ALTER TABLE "identity_verifiable_addresses" ADD COLUMN "expires_at" DATETIME;
UPDATE identity_verifiable_addresses SET code = substr(hex(randomblob(32)), 0, 32) WHERE code IS NULL;
UPDATE identity_verifiable_addresses SET expires_at = CURRENT_TIMESTAMP WHERE expires_at IS NULL;
DROP INDEX IF EXISTS "identity_verifiable_addresses_status_via_uq_idx";
DROP INDEX IF EXISTS "identity_verifiable_addresses_status_via_idx";
CREATE TABLE "_identity_verifiable_addresses_tmp" (
"id" TEXT PRIMARY KEY,
"status" TEXT NOT NULL,
"via" TEXT NOT NULL,
"verified" bool NOT NULL,
"value" TEXT NOT NULL,
"verified_at" DATETIME,
"identity_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"code" TEXT NOT NULL,
"expires_at" DATETIME,
FOREIGN KEY (identity_id) REFERENCES identities (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE UNIQUE INDEX "identity_verifiable_addresses_status_via_uq_idx" ON "_identity_verifiable_addresses_tmp" (via, value);
CREATE INDEX "identity_verifiable_addresses_status_via_idx" ON "_identity_verifiable_addresses_tmp" (via, value);
INSERT INTO "_identity_verifiable_addresses_tmp" (id, status, via, verified, value, verified_at, identity_id, created_at, updated_at, code, expires_at) SELECT id, status, via, verified, value, verified_at, identity_id, created_at, updated_at, code, expires_at FROM "identity_verifiable_addresses";
DROP TABLE "identity_verifiable_addresses";
ALTER TABLE "_identity_verifiable_addresses_tmp" RENAME TO "identity_verifiable_addresses";
DROP INDEX IF EXISTS "identity_verifiable_addresses_status_via_uq_idx";
DROP INDEX IF EXISTS "identity_verifiable_addresses_status_via_idx";
CREATE TABLE "_identity_verifiable_addresses_tmp" (
"id" TEXT PRIMARY KEY,
"status" TEXT NOT NULL,
"via" TEXT NOT NULL,
"verified" bool NOT NULL,
"value" TEXT NOT NULL,
"verified_at" DATETIME,
"identity_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"code" TEXT NOT NULL,
"expires_at" DATETIME,
FOREIGN KEY (identity_id) REFERENCES identities (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE UNIQUE INDEX "identity_verifiable_addresses_status_via_uq_idx" ON "_identity_verifiable_addresses_tmp" (via, value);
CREATE INDEX "identity_verifiable_addresses_status_via_idx" ON "_identity_verifiable_addresses_tmp" (via, value);
INSERT INTO "_identity_verifiable_addresses_tmp" (id, status, via, verified, value, verified_at, identity_id, created_at, updated_at, code, expires_at) SELECT id, status, via, verified, value, verified_at, identity_id, created_at, updated_at, code, expires_at FROM "identity_verifiable_addresses";
DROP TABLE "identity_verifiable_addresses";
ALTER TABLE "_identity_verifiable_addresses_tmp" RENAME TO "identity_verifiable_addresses";
CREATE UNIQUE INDEX "identity_verifiable_addresses_code_uq_idx" ON "identity_verifiable_addresses" (code);
CREATE INDEX "identity_verifiable_addresses_code_idx" ON "identity_verifiable_addresses" (code);