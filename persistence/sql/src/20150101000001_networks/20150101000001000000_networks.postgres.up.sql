CREATE TABLE "networks" (
  "id" UUID NOT NULL,
  PRIMARY KEY("id"),
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL
);

INSERT INTO networks (id, created_at, updated_at) VALUES (uuid_in(
  overlay(
    overlay(
      md5(random()::text || ':' || clock_timestamp()::text)
      placing '4'
      from 13
    )
    placing to_hex(floor(random()*(11-8+1) + 8)::int)::text
    from 17
  )::cstring
), '2013-10-07 08:23:19', '2013-10-07 08:23:19');
