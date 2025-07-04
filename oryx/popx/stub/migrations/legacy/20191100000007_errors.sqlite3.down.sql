CREATE TABLE "_selfservice_errors_tmp" (
"id" TEXT PRIMARY KEY,
"errors" TEXT NOT NULL,
"seen_at" DATETIME,
"was_seen" bool NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
INSERT INTO "_selfservice_errors_tmp" (id, errors, seen_at, was_seen, created_at, updated_at) SELECT id, errors, seen_at, was_seen, created_at, updated_at FROM "selfservice_errors";

DROP TABLE "selfservice_errors";
ALTER TABLE "_selfservice_errors_tmp" RENAME TO "selfservice_errors";