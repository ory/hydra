CREATE TABLE "selfservice_recovery_requests" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"request_url" VARCHAR (2048) NOT NULL,
"issued_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
"expires_at" timestamp NOT NULL,
"messages" jsonb,
"active_method" VARCHAR (32),
"csrf_token" VARCHAR (255) NOT NULL,
"state" VARCHAR (32) NOT NULL,
"recovered_identity_id" UUID,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
FOREIGN KEY ("recovered_identity_id") REFERENCES "identities" ("id") ON DELETE cascade
)