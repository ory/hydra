CREATE TABLE "continuity_containers" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"identity_id" UUID,
"name" VARCHAR (255) NOT NULL,
"payload" json,
"expires_at" timestamp NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
CONSTRAINT "continuity_containers_identities_id_fk" FOREIGN KEY ("identity_id") REFERENCES "identities" ("id") ON DELETE cascade
);COMMIT TRANSACTION;BEGIN TRANSACTION;