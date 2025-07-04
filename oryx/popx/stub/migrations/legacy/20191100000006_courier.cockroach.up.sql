CREATE TABLE "courier_messages" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"type" int NOT NULL,
"status" int NOT NULL,
"body" VARCHAR (255) NOT NULL,
"subject" VARCHAR (255) NOT NULL,
"recipient" VARCHAR (255) NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL
);COMMIT TRANSACTION;BEGIN TRANSACTION;