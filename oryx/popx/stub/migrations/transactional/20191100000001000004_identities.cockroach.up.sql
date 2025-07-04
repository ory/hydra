CREATE TABLE "identity_credential_identifiers" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"identifier" VARCHAR (255) NOT NULL,
"identity_credential_id" UUID NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
CONSTRAINT "identity_credential_identifiers_identity_credentials_id_fk" FOREIGN KEY ("identity_credential_id") REFERENCES "identity_credentials" ("id") ON DELETE cascade
)