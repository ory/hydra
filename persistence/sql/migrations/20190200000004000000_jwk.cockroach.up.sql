CREATE TABLE IF NOT EXISTS hydra_jwk (
	sid varchar(255) NOT NULL,
	kid varchar(255) NOT NULL,
	version int NOT NULL DEFAULT 0,
	keydata text NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	pk SERIAL,
	UNIQUE INDEX (sid, kid),
  CONSTRAINT "primary" PRIMARY KEY (pk ASC)
);
