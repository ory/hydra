CREATE TABLE IF NOT EXISTS hydra_oauth2_jti_blacklist (
	signature       varchar(64) NOT NULL PRIMARY KEY,
	expires_at  	timestamp NOT NULL DEFAULT now()
);

-- mysql requires the index to be named
CREATE INDEX hydra_oauth2_jti_blacklist_expiry ON hydra_oauth2_jti_blacklist ( expires_at );

