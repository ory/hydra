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
)