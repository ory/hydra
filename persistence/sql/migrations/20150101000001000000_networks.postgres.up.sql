-- Migration generated by the command below; DO NOT EDIT.
-- hydra:generate hydra migrate gen
CREATE TABLE "networks" (
  "id" UUID NOT NULL,
  PRIMARY KEY("id"),
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL
);
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
INSERT INTO networks (id, created_at, updated_at) VALUES (uuid_generate_v4(), '2013-10-07 08:23:19', '2013-10-07 08:23:19');
