CREATE TABLE "selfservice_login_requests" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"request_url" VARCHAR (2048) NOT NULL,
"issued_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
"expires_at" timestamp NOT NULL,
"active_method" VARCHAR (32) NOT NULL,
"csrf_token" VARCHAR (255) NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL
);
CREATE TABLE "selfservice_login_request_methods" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"method" VARCHAR (32) NOT NULL,
"selfservice_login_request_id" UUID NOT NULL,
"config" jsonb NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
FOREIGN KEY ("selfservice_login_request_id") REFERENCES "selfservice_login_requests" ("id") ON DELETE cascade
);
CREATE TABLE "selfservice_registration_requests" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"request_url" VARCHAR (2048) NOT NULL,
"issued_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
"expires_at" timestamp NOT NULL,
"active_method" VARCHAR (32) NOT NULL,
"csrf_token" VARCHAR (255) NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL
);
CREATE TABLE "selfservice_registration_request_methods" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"method" VARCHAR (32) NOT NULL,
"selfservice_registration_request_id" UUID NOT NULL,
"config" jsonb NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
FOREIGN KEY ("selfservice_registration_request_id") REFERENCES "selfservice_registration_requests" ("id") ON DELETE cascade
);
CREATE TABLE "selfservice_profile_management_requests" (
"id" UUID NOT NULL,
PRIMARY KEY("id"),
"request_url" VARCHAR (2048) NOT NULL,
"issued_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
"expires_at" timestamp NOT NULL,
"form" jsonb NOT NULL,
"update_successful" bool NOT NULL,
"identity_id" UUID NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
FOREIGN KEY ("identity_id") REFERENCES "identities" ("id") ON DELETE cascade
);