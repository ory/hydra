CREATE TABLE "identity_verification_tokens" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"token" VARCHAR (64) NOT NULL,
"used" bool NOT NULL DEFAULT 'false',
"used_at" timestamp,
"expires_at" timestamp NOT NULL,
"issued_at" timestamp NOT NULL,
"identity_verifiable_address_id" UUID NOT NULL,
"selfservice_verification_flow_id" UUID,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
CONSTRAINT "identity_verification_tokens_identity_verifiable_addresses_id_fk" FOREIGN KEY ("identity_verifiable_address_id") REFERENCES "identity_verifiable_addresses" ("id") ON DELETE cascade,
CONSTRAINT "identity_verification_tokens_selfservice_verification_flows_id_fk" FOREIGN KEY ("selfservice_verification_flow_id") REFERENCES "selfservice_verification_flows" ("id") ON DELETE cascade
);COMMIT TRANSACTION;BEGIN TRANSACTION;
CREATE UNIQUE INDEX "identity_verification_tokens_token_uq_idx" ON "identity_verification_tokens" (token);COMMIT TRANSACTION;BEGIN TRANSACTION;
CREATE INDEX "identity_verification_tokens_token_idx" ON "identity_verification_tokens" (token);COMMIT TRANSACTION;BEGIN TRANSACTION;
CREATE INDEX "identity_verification_tokens_verifiable_address_id_idx" ON "identity_verification_tokens" (identity_verifiable_address_id);COMMIT TRANSACTION;BEGIN TRANSACTION;
CREATE INDEX "identity_verification_tokens_verification_flow_id_idx" ON "identity_verification_tokens" (selfservice_verification_flow_id);COMMIT TRANSACTION;BEGIN TRANSACTION;