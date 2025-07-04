CREATE TABLE "identity_recovery_addresses" (
"id" TEXT PRIMARY KEY,
"via" TEXT NOT NULL,
"value" TEXT NOT NULL,
"identity_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (identity_id) REFERENCES identities (id) ON DELETE cascade
)