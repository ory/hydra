CREATE TABLE "selfservice_login_request_methods" (
"id" TEXT PRIMARY KEY,
"method" TEXT NOT NULL,
"selfservice_login_request_id" char(36) NOT NULL,
"config" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (selfservice_login_request_id) REFERENCES selfservice_login_requests (id) ON DELETE cascade
)