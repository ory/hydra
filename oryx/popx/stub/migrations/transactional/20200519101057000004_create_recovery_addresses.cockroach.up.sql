CREATE TABLE "selfservice_recovery_request_methods" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"method" VARCHAR (32) NOT NULL,
"config" json NOT NULL,
"selfservice_recovery_request_id" UUID NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
CONSTRAINT "selfservice_recovery_request_methods_selfservice_recovery_requests_id_fk" FOREIGN KEY ("selfservice_recovery_request_id") REFERENCES "selfservice_recovery_requests" ("id") ON DELETE cascade
)