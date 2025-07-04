CREATE TABLE "selfservice_profile_management_requests" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"request_url" VARCHAR (2048) NOT NULL,
"issued_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
"expires_at" timestamp NOT NULL,
"form" json NOT NULL,
"update_successful" bool NOT NULL,
"identity_id" UUID NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
CONSTRAINT "selfservice_profile_management_requests_identities_id_fk" FOREIGN KEY ("identity_id") REFERENCES "identities" ("id") ON DELETE cascade
);