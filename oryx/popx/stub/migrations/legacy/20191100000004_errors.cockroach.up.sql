CREATE TABLE "selfservice_errors" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"errors" json NOT NULL,
"seen_at" timestamp NOT NULL,
"was_seen" bool NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL
);COMMIT TRANSACTION;BEGIN TRANSACTION;