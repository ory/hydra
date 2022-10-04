CREATE TABLE "networks" (
  "id" UUID NOT NULL,
  PRIMARY KEY("id"),
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL
);

INSERT INTO networks (id, created_at, updated_at) VALUES (gen_random_uuid(), '2013-10-07 08:23:19', '2013-10-07 08:23:19');
