CREATE TABLE IF NOT EXISTS hydra_oauth2_access (
	signature varchar(255) NOT NULL,
	request_id varchar(40) NOT NULL,
	requested_at timestamp NOT NULL DEFAULT now(),
	client_id varchar(255) NOT NULL,
	scope text NOT NULL,
	granted_scope text NOT NULL,
	form_data text NOT NULL,
	session_data text NOT NULL,
	subject varchar(255) NOT NULL DEFAULT '',
	active bool NOT NULL DEFAULT TRUE,
	requested_audience text NULL DEFAULT '',
	granted_audience text NULL DEFAULT '',
	challenge_id varchar(40) NULL,
	UNIQUE (request_id),
	INDEX (requested_at),
	INDEX (client_id),
	INDEX (challenge_id),
  CONSTRAINT "primary" PRIMARY KEY (signature ASC)
);
CREATE TABLE IF NOT EXISTS hydra_oauth2_refresh (
	signature varchar(255) NOT NULL,
	request_id varchar(40) NOT NULL,
	requested_at timestamp NOT NULL DEFAULT now(),
	client_id varchar(255) NOT NULL,
	scope text NOT NULL,
	granted_scope text NOT NULL,
	form_data text NOT NULL,
	session_data text NOT NULL,
	subject varchar(255) NOT NULL DEFAULT '',
	active bool NOT NULL DEFAULT TRUE,
	requested_audience text NULL DEFAULT '',
	granted_audience text NULL DEFAULT '',
	challenge_id varchar(40) NULL,
	UNIQUE (request_id),
	INDEX (client_id),
	INDEX (challenge_id),
  CONSTRAINT "primary" PRIMARY KEY (signature ASC)
);
CREATE TABLE IF NOT EXISTS hydra_oauth2_code (
	signature varchar(255) NOT NULL,
	request_id varchar(40) NOT NULL,
	requested_at timestamp NOT NULL DEFAULT now(),
	client_id varchar(255) NOT NULL,
	scope text NOT NULL,
	granted_scope text NOT NULL,
	form_data text NOT NULL,
	session_data text NOT NULL,
	subject varchar(255) NOT NULL DEFAULT '',
	active bool NOT NULL DEFAULT TRUE,
	requested_audience text NULL DEFAULT '',
	granted_audience text NULL DEFAULT '',
	challenge_id varchar(40) NULL,
	INDEX (client_id),
	INDEX (challenge_id),
  CONSTRAINT "primary" PRIMARY KEY (signature ASC)
);
CREATE TABLE IF NOT EXISTS hydra_oauth2_oidc (
	signature varchar(255) NOT NULL,
	request_id varchar(40) NOT NULL,
	requested_at timestamp NOT NULL DEFAULT now(),
	client_id varchar(255) NOT NULL,
	scope text NOT NULL,
	granted_scope text NOT NULL,
	form_data text NOT NULL,
	session_data text NOT NULL,
	subject varchar(255) NOT NULL DEFAULT '',
	active bool NOT NULL DEFAULT TRUE,
	requested_audience text NULL DEFAULT '',
	granted_audience text NULL DEFAULT '',
	challenge_id varchar(40) NULL,
	INDEX (client_id),
	INDEX (challenge_id),
  CONSTRAINT "primary" PRIMARY KEY (signature ASC)
);
CREATE TABLE IF NOT EXISTS hydra_oauth2_pkce (
	signature varchar(255) NOT NULL,
	request_id varchar(40) NOT NULL,
	requested_at timestamp NOT NULL DEFAULT now(),
	client_id varchar(255) NOT NULL,
	scope text NOT NULL,
	granted_scope text NOT NULL,
	form_data text NOT NULL,
	session_data text NOT NULL,
	subject varchar(255) NOT NULL,
	active bool NOT NULL DEFAULT TRUE,
	requested_audience TEXT NULL DEFAULT '',
	granted_audience TEXT NULL DEFAULT '',
	challenge_id varchar(40) NULL,
	INDEX (client_id),
	INDEX (challenge_id),
  CONSTRAINT "primary" PRIMARY KEY (signature ASC)
);

ALTER TABLE hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;

ALTER TABLE hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
