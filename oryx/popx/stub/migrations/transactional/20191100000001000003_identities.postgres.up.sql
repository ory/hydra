CREATE TABLE "identity_credentials" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"config" jsonb NOT NULL,
"identity_credential_type_id" UUID NOT NULL,
"identity_id" UUID NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
FOREIGN KEY ("identity_id") REFERENCES "identities" ("id") ON DELETE cascade,
FOREIGN KEY ("identity_credential_type_id") REFERENCES "identity_credential_types" ("id") ON DELETE cascade
)