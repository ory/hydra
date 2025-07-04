CREATE TABLE "identities" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"traits_schema_id" VARCHAR (2048) NOT NULL,
"traits" jsonb NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL
)