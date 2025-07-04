CREATE TABLE "identity_verifiable_addresses" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"code" VARCHAR (32) NOT NULL,
"status" VARCHAR (16) NOT NULL,
"via" VARCHAR (16) NOT NULL,
"verified" bool NOT NULL,
"value" VARCHAR (400) NOT NULL,
"verified_at" timestamp,
"expires_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
"identity_id" UUID NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
FOREIGN KEY ("identity_id") REFERENCES "identities" ("id") ON DELETE cascade
)