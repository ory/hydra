CREATE TABLE "selfservice_verification_requests" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"request_url" VARCHAR (2048) NOT NULL,
"issued_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
"expires_at" timestamp NOT NULL,
"form" jsonb NOT NULL,
"via" VARCHAR (16) NOT NULL,
"csrf_token" VARCHAR (255) NOT NULL,
"success" bool NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL
);