CREATE TABLE "_selfservice_login_requests_tmp" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"active_method" TEXT NOT NULL,
"csrf_token" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
INSERT INTO "_selfservice_login_requests_tmp" (id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at) SELECT id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at FROM "selfservice_login_requests";

DROP TABLE "selfservice_login_requests";
ALTER TABLE "_selfservice_login_requests_tmp" RENAME TO "selfservice_login_requests";