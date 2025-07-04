CREATE TABLE "selfservice_registration_request_methods" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"method" VARCHAR (32) NOT NULL,
"selfservice_registration_request_id" UUID NOT NULL,
"config" jsonb NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
FOREIGN KEY ("selfservice_registration_request_id") REFERENCES "selfservice_registration_requests" ("id") ON DELETE cascade
)