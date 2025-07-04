CREATE TABLE "selfservice_verification_flow_methods" (
"id" TEXT PRIMARY KEY,
"method" TEXT NOT NULL,
"selfservice_verification_flow_id" char(36) NOT NULL,
"config" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
ALTER TABLE "selfservice_verification_flows" ADD COLUMN "active_method" TEXT;