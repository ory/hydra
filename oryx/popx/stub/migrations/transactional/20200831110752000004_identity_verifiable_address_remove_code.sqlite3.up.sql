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
)