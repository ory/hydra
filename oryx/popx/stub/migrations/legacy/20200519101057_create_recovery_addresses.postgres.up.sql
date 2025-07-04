CREATE TABLE "identity_recovery_addresses" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"via" VARCHAR (16) NOT NULL,
"value" VARCHAR (400) NOT NULL,
"identity_id" UUID NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
FOREIGN KEY ("identity_id") REFERENCES "identities" ("id") ON DELETE cascade
);
CREATE UNIQUE INDEX "identity_recovery_addresses_status_via_uq_idx" ON "identity_recovery_addresses" (via, value);
CREATE INDEX "identity_recovery_addresses_status_via_idx" ON "identity_recovery_addresses" (via, value);
CREATE TABLE "selfservice_recovery_requests" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"request_url" VARCHAR (2048) NOT NULL,
"issued_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
"expires_at" timestamp NOT NULL,
"messages" jsonb,
"active_method" VARCHAR (32),
"csrf_token" VARCHAR (255) NOT NULL,
"state" VARCHAR (32) NOT NULL,
"recovered_identity_id" UUID,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
FOREIGN KEY ("recovered_identity_id") REFERENCES "identities" ("id") ON DELETE cascade
);
CREATE TABLE "selfservice_recovery_request_methods" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"method" VARCHAR (32) NOT NULL,
"config" jsonb NOT NULL,
"selfservice_recovery_request_id" UUID NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
FOREIGN KEY ("selfservice_recovery_request_id") REFERENCES "selfservice_recovery_requests" ("id") ON DELETE cascade
);
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
FOREIGN KEY ("identity_recovery_address_id") REFERENCES "identity_recovery_addresses" ("id") ON DELETE cascade,
FOREIGN KEY ("selfservice_recovery_request_id") REFERENCES "selfservice_recovery_requests" ("id") ON DELETE cascade
);
CREATE UNIQUE INDEX "identity_recovery_addresses_code_uq_idx" ON "identity_recovery_tokens" (token);
CREATE INDEX "identity_recovery_addresses_code_idx" ON "identity_recovery_tokens" (token);