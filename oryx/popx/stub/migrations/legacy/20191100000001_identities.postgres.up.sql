CREATE TABLE "identities" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"traits_schema_id" VARCHAR (2048) NOT NULL,
"traits" jsonb NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL
);
CREATE TABLE "identity_credential_types" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"name" VARCHAR (32) NOT NULL
);
CREATE UNIQUE INDEX "identity_credential_types_name_idx" ON "identity_credential_types" (name);
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
);
CREATE TABLE "identity_credential_identifiers" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"identifier" VARCHAR (255) NOT NULL,
"identity_credential_id" UUID NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
FOREIGN KEY ("identity_credential_id") REFERENCES "identity_credentials" ("id") ON DELETE cascade
);
CREATE UNIQUE INDEX "identity_credential_identifiers_identifier_idx" ON "identity_credential_identifiers" (identifier);