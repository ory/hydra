CREATE TABLE "courier_messages" (
"id" TEXT PRIMARY KEY,
"type" INTEGER NOT NULL,
"status" INTEGER NOT NULL,
"body" TEXT NOT NULL,
"subject" TEXT NOT NULL,
"recipient" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);