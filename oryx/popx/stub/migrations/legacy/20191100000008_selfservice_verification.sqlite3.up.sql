CREATE TABLE "identity_verifiable_addresses" (
"id" TEXT PRIMARY KEY,
"code" TEXT NOT NULL,
"status" TEXT NOT NULL,
"via" TEXT NOT NULL,
"verified" bool NOT NULL,
"value" TEXT NOT NULL,
"verified_at" DATETIME,
"expires_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
"identity_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (identity_id) REFERENCES identities (id) ON DELETE cascade
);
CREATE UNIQUE INDEX "identity_verifiable_addresses_code_uq_idx" ON "identity_verifiable_addresses" (code);
CREATE INDEX "identity_verifiable_addresses_code_idx" ON "identity_verifiable_addresses" (code);
CREATE UNIQUE INDEX "identity_verifiable_addresses_status_via_uq_idx" ON "identity_verifiable_addresses" (via, value);
CREATE INDEX "identity_verifiable_addresses_status_via_idx" ON "identity_verifiable_addresses" (via, value);
CREATE TABLE "selfservice_verification_requests" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
"expires_at" DATETIME NOT NULL,
"form" TEXT NOT NULL,
"via" TEXT NOT NULL,
"csrf_token" TEXT NOT NULL,
"success" bool NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);