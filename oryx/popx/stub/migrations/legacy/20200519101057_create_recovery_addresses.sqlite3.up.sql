CREATE TABLE "identity_recovery_addresses" (
"id" TEXT PRIMARY KEY,
"via" TEXT NOT NULL,
"value" TEXT NOT NULL,
"identity_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (identity_id) REFERENCES identities (id) ON DELETE cascade
);
CREATE UNIQUE INDEX "identity_recovery_addresses_status_via_uq_idx" ON "identity_recovery_addresses" (via, value);
CREATE INDEX "identity_recovery_addresses_status_via_idx" ON "identity_recovery_addresses" (via, value);
CREATE TABLE "selfservice_recovery_requests" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
"expires_at" DATETIME NOT NULL,
"messages" TEXT,
"active_method" TEXT,
"csrf_token" TEXT NOT NULL,
"state" TEXT NOT NULL,
"recovered_identity_id" char(36),
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (recovered_identity_id) REFERENCES identities (id) ON DELETE cascade
);
CREATE TABLE "selfservice_recovery_request_methods" (
"id" TEXT PRIMARY KEY,
"method" TEXT NOT NULL,
"config" TEXT NOT NULL,
"selfservice_recovery_request_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (selfservice_recovery_request_id) REFERENCES selfservice_recovery_requests (id) ON DELETE cascade
);
CREATE TABLE "identity_recovery_tokens" (
"id" TEXT PRIMARY KEY,
"token" TEXT NOT NULL,
"used" bool NOT NULL DEFAULT 'false',
"used_at" DATETIME,
"identity_recovery_address_id" char(36) NOT NULL,
"selfservice_recovery_request_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (identity_recovery_address_id) REFERENCES identity_recovery_addresses (id) ON DELETE cascade,
FOREIGN KEY (selfservice_recovery_request_id) REFERENCES selfservice_recovery_requests (id) ON DELETE cascade
);
CREATE UNIQUE INDEX "identity_recovery_addresses_code_uq_idx" ON "identity_recovery_tokens" (token);
CREATE INDEX "identity_recovery_addresses_code_idx" ON "identity_recovery_tokens" (token);