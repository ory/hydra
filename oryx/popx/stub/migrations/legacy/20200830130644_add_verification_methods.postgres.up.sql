CREATE TABLE "selfservice_verification_flow_methods" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"method" VARCHAR (32) NOT NULL,
"selfservice_verification_flow_id" UUID NOT NULL,
"config" jsonb NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL
);
ALTER TABLE "selfservice_verification_flows" ADD COLUMN "active_method" VARCHAR (32);