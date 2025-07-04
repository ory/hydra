CREATE TABLE "selfservice_recovery_request_methods" (
"id" TEXT PRIMARY KEY,
"method" TEXT NOT NULL,
"config" TEXT NOT NULL,
"selfservice_recovery_request_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (selfservice_recovery_request_id) REFERENCES selfservice_recovery_requests (id) ON DELETE cascade
)