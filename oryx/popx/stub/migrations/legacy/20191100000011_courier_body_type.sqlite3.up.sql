CREATE TABLE "_courier_messages_tmp" (
"id" TEXT PRIMARY KEY,
"type" INTEGER NOT NULL,
"status" INTEGER NOT NULL,
"body" TEXT NOT NULL,
"subject" TEXT NOT NULL,
"recipient" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
INSERT INTO "_courier_messages_tmp" (id, type, status, body, subject, recipient, created_at, updated_at) SELECT id, type, status, body, subject, recipient, created_at, updated_at FROM "courier_messages";
DROP TABLE "courier_messages";
ALTER TABLE "_courier_messages_tmp" RENAME TO "courier_messages";