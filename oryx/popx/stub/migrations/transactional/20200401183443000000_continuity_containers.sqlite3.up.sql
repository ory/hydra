CREATE TABLE "continuity_containers" (
"id" TEXT PRIMARY KEY,
"identity_id" char(36),
"name" TEXT NOT NULL,
"payload" TEXT,
"expires_at" DATETIME NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (identity_id) REFERENCES identities (id) ON DELETE cascade
);