DROP TABLE "selfservice_verification_flow_methods";
CREATE TABLE "_selfservice_verification_flows_tmp" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"csrf_token" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"messages" TEXT,
"type" TEXT NOT NULL DEFAULT 'browser',
"state" TEXT NOT NULL DEFAULT 'show_form'
);
INSERT INTO "_selfservice_verification_flows_tmp" (id, request_url, issued_at, expires_at, csrf_token, created_at, updated_at, messages, type, state) SELECT id, request_url, issued_at, expires_at, csrf_token, created_at, updated_at, messages, type, state FROM "selfservice_verification_flows";

DROP TABLE "selfservice_verification_flows";
ALTER TABLE "_selfservice_verification_flows_tmp" RENAME TO "selfservice_verification_flows";
CREATE TABLE "_selfservice_verification_flows_tmp" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"csrf_token" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"messages" TEXT,
"type" TEXT NOT NULL DEFAULT 'browser'
);
INSERT INTO "_selfservice_verification_flows_tmp" (id, request_url, issued_at, expires_at, csrf_token, created_at, updated_at, messages, type) SELECT id, request_url, issued_at, expires_at, csrf_token, created_at, updated_at, messages, type FROM "selfservice_verification_flows";

DROP TABLE "selfservice_verification_flows";
ALTER TABLE "_selfservice_verification_flows_tmp" RENAME TO "selfservice_verification_flows";
ALTER TABLE "selfservice_verification_flows" ADD COLUMN "via" TEXT NOT NULL DEFAULT 'email';
ALTER TABLE "selfservice_verification_flows" ADD COLUMN "success" bool NOT NULL DEFAULT FALSE;