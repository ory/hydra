CREATE TABLE "_selfservice_settings_requests_tmp" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"identity_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"active_method" TEXT,
"update_successful" bool NOT NULL DEFAULT 'false',
FOREIGN KEY (identity_id) REFERENCES identities (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
INSERT INTO "_selfservice_settings_requests_tmp" (id, request_url, issued_at, expires_at, identity_id, created_at, updated_at, active_method, update_successful) SELECT id, request_url, issued_at, expires_at, identity_id, created_at, updated_at, active_method, update_successful FROM "selfservice_settings_requests";

DROP TABLE "selfservice_settings_requests";
ALTER TABLE "_selfservice_settings_requests_tmp" RENAME TO "selfservice_settings_requests";