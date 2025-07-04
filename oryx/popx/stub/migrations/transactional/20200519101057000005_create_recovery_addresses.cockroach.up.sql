CREATE TABLE "identity_recovery_tokens" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"token" VARCHAR (64) NOT NULL,
"used" bool NOT NULL DEFAULT 'false',
"used_at" timestamp,
"identity_recovery_address_id" UUID NOT NULL,
"selfservice_recovery_request_id" UUID NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
CONSTRAINT "identity_recovery_tokens_identity_recovery_addresses_id_fk" FOREIGN KEY ("identity_recovery_address_id") REFERENCES "identity_recovery_addresses" ("id") ON DELETE cascade,
CONSTRAINT "identity_recovery_tokens_selfservice_recovery_requests_id_fk" FOREIGN KEY ("selfservice_recovery_request_id") REFERENCES "selfservice_recovery_requests" ("id") ON DELETE cascade
)