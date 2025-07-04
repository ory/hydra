DELETE FROM sessions;
ALTER TABLE "sessions" ADD COLUMN "token" TEXT;
CREATE TABLE "_sessions_tmp" (
"id" TEXT PRIMARY KEY,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"authenticated_at" DATETIME NOT NULL,
"identity_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"token" TEXT,
FOREIGN KEY (identity_id) REFERENCES identities (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
INSERT INTO "_sessions_tmp" (id, issued_at, expires_at, authenticated_at, identity_id, created_at, updated_at, token) SELECT id, issued_at, expires_at, authenticated_at, identity_id, created_at, updated_at, token FROM "sessions";
DROP TABLE "sessions";
ALTER TABLE "_sessions_tmp" RENAME TO "sessions";
CREATE UNIQUE INDEX "sessions_token_uq_idx" ON "sessions" (token);
CREATE INDEX "sessions_token_idx" ON "sessions" (token);