DROP INDEX IF EXISTS "identity_verifiable_addresses_code_uq_idx";
DROP INDEX IF EXISTS "identity_verifiable_addresses_code_idx";
DROP INDEX IF EXISTS "identity_verifiable_addresses_status_via_idx";
DROP INDEX IF EXISTS "identity_verifiable_addresses_status_via_uq_idx";
CREATE TABLE "_identity_verifiable_addresses_tmp" (
"id" TEXT PRIMARY KEY,
"status" TEXT NOT NULL,
"via" TEXT NOT NULL,
"verified" bool NOT NULL,
"value" TEXT NOT NULL,
"verified_at" DATETIME,
"expires_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"identity_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (identity_id) REFERENCES identities (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "identity_verifiable_addresses_status_via_idx" ON "_identity_verifiable_addresses_tmp" (via, value);
CREATE UNIQUE INDEX "identity_verifiable_addresses_status_via_uq_idx" ON "_identity_verifiable_addresses_tmp" (via, value);
INSERT INTO "_identity_verifiable_addresses_tmp" (id, status, via, verified, value, verified_at, expires_at, identity_id, created_at, updated_at) SELECT id, status, via, verified, value, verified_at, expires_at, identity_id, created_at, updated_at FROM "identity_verifiable_addresses";

DROP TABLE "identity_verifiable_addresses";
ALTER TABLE "_identity_verifiable_addresses_tmp" RENAME TO "identity_verifiable_addresses";
DROP INDEX IF EXISTS "identity_verifiable_addresses_status_via_idx";
DROP INDEX IF EXISTS "identity_verifiable_addresses_status_via_uq_idx";
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
FOREIGN KEY (identity_id) REFERENCES identities (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "identity_verifiable_addresses_status_via_idx" ON "_identity_verifiable_addresses_tmp" (via, value);
CREATE UNIQUE INDEX "identity_verifiable_addresses_status_via_uq_idx" ON "_identity_verifiable_addresses_tmp" (via, value);
INSERT INTO "_identity_verifiable_addresses_tmp" (id, status, via, verified, value, verified_at, identity_id, created_at, updated_at) SELECT id, status, via, verified, value, verified_at, identity_id, created_at, updated_at FROM "identity_verifiable_addresses";

DROP TABLE "identity_verifiable_addresses";
ALTER TABLE "_identity_verifiable_addresses_tmp" RENAME TO "identity_verifiable_addresses";