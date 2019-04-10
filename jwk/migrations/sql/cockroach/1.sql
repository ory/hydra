-- +migrate Up
CREATE TABLE IF NOT EXISTS hydra_jwk (
  pk SERIAL PRIMARY KEY,
	sid varchar(255) NOT NULL,
	kid varchar(255) NOT NULL,
	version int NOT NULL DEFAULT 0,
	keydata text NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	UNIQUE INDEX (sid, kid)
);

-- +migrate Down
DROP TABLE hydra_jwk;
