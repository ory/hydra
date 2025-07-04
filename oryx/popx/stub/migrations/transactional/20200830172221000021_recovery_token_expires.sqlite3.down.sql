CREATE TABLE "_identity_recovery_tokens_tmp" (
"id" TEXT PRIMARY KEY,
"token" TEXT NOT NULL,
"used" bool NOT NULL DEFAULT 'false',
"used_at" DATETIME,
"identity_recovery_address_id" char(36) NOT NULL,
"selfservice_recovery_flow_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"expires_at" DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
"issued_at" DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
FOREIGN KEY (identity_recovery_address_id) REFERENCES identity_recovery_addresses (id) ON UPDATE NO ACTION ON DELETE CASCADE,
FOREIGN KEY (selfservice_recovery_flow_id) REFERENCES selfservice_recovery_flows (id) ON UPDATE NO ACTION ON DELETE CASCADE
)