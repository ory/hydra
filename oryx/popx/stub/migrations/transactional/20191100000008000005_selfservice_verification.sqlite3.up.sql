CREATE TABLE "selfservice_verification_requests" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
"expires_at" DATETIME NOT NULL,
"form" TEXT NOT NULL,
"via" TEXT NOT NULL,
"csrf_token" TEXT NOT NULL,
"success" bool NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);