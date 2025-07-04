CREATE TABLE "_selfservice_verification_requests_tmp" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"csrf_token" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"via" TEXT NOT NULL DEFAULT 'email',
"success" bool NOT NULL DEFAULT 'FALSE'
);
INSERT INTO "_selfservice_verification_requests_tmp" (id, request_url, issued_at, expires_at, csrf_token, created_at, updated_at, via, success) SELECT id, request_url, issued_at, expires_at, csrf_token, created_at, updated_at, via, success FROM "selfservice_verification_requests";

DROP TABLE "selfservice_verification_requests";
ALTER TABLE "_selfservice_verification_requests_tmp" RENAME TO "selfservice_verification_requests";
CREATE TABLE "_selfservice_login_requests_tmp" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"active_method" TEXT NOT NULL,
"csrf_token" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"forced" bool NOT NULL DEFAULT 'false'
);
INSERT INTO "_selfservice_login_requests_tmp" (id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at, forced) SELECT id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at, forced FROM "selfservice_login_requests";

DROP TABLE "selfservice_login_requests";
ALTER TABLE "_selfservice_login_requests_tmp" RENAME TO "selfservice_login_requests";
CREATE TABLE "_selfservice_registration_requests_tmp" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"active_method" TEXT NOT NULL,
"csrf_token" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
INSERT INTO "_selfservice_registration_requests_tmp" (id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at) SELECT id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at FROM "selfservice_registration_requests";

DROP TABLE "selfservice_registration_requests";
ALTER TABLE "_selfservice_registration_requests_tmp" RENAME TO "selfservice_registration_requests";