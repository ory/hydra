CREATE TABLE "identity_credentials" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"config" json NOT NULL,
"identity_credential_type_id" UUID NOT NULL,
"identity_id" UUID NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
CONSTRAINT "identity_credentials_identities_id_fk" FOREIGN KEY ("identity_id") REFERENCES "identities" ("id") ON DELETE cascade,
CONSTRAINT "identity_credentials_identity_credential_types_id_fk" FOREIGN KEY ("identity_credential_type_id") REFERENCES "identity_credential_types" ("id") ON DELETE cascade
)