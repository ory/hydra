CREATE TABLE "identities" (
"id" TEXT PRIMARY KEY,
"traits_schema_id" TEXT NOT NULL,
"traits" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
)