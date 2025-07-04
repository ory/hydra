CREATE TABLE "identities" (
"id" TEXT PRIMARY KEY,
"traits_schema_id" TEXT NOT NULL,
"traits" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE "identity_credential_types" (
"id" TEXT PRIMARY KEY,
"name" TEXT NOT NULL
);
CREATE UNIQUE INDEX "identity_credential_types_name_idx" ON "identity_credential_types" (name);
CREATE TABLE "identity_credentials" (
"id" TEXT PRIMARY KEY,
"config" TEXT NOT NULL,
"identity_credential_type_id" char(36) NOT NULL,
"identity_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (identity_id) REFERENCES identities (id) ON DELETE cascade,
FOREIGN KEY (identity_credential_type_id) REFERENCES identity_credential_types (id) ON DELETE cascade
);
CREATE TABLE "identity_credential_identifiers" (
"id" TEXT PRIMARY KEY,
"identifier" TEXT NOT NULL,
"identity_credential_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (identity_credential_id) REFERENCES identity_credentials (id) ON DELETE cascade
);
CREATE UNIQUE INDEX "identity_credential_identifiers_identifier_idx" ON "identity_credential_identifiers" (identifier);