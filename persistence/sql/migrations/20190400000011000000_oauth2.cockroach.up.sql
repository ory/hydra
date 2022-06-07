CREATE TABLE IF NOT EXISTS hydra_oauth2_jti_blacklist (
	signature      	varchar(64) NOT NULL,
	expires_at  	timestamp NOT NULL DEFAULT now(),
  CONSTRAINT "primary" PRIMARY KEY (signature ASC)
);

CREATE INDEX ON hydra_oauth2_jti_blacklist ( expires_at );
