CREATE TABLE "identity_recovery_addresses" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"via" VARCHAR (16) NOT NULL,
"value" VARCHAR (400) NOT NULL,
"identity_id" UUID NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
CONSTRAINT "identity_recovery_addresses_identities_id_fk" FOREIGN KEY ("identity_id") REFERENCES "identities" ("id") ON DELETE cascade
)