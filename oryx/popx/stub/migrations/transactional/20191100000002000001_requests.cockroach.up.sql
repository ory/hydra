CREATE TABLE "selfservice_login_request_methods" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"method" VARCHAR (32) NOT NULL,
"selfservice_login_request_id" UUID NOT NULL,
"config" json NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
CONSTRAINT "selfservice_login_request_methods_selfservice_login_requests_id_fk" FOREIGN KEY ("selfservice_login_request_id") REFERENCES "selfservice_login_requests" ("id") ON DELETE cascade
)