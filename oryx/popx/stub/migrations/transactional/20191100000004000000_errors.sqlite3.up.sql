CREATE TABLE "selfservice_errors" (
"id" TEXT PRIMARY KEY,
"errors" TEXT NOT NULL,
"seen_at" DATETIME NOT NULL,
"was_seen" bool NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);