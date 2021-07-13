CREATE TABLE hydra_oauth2_consent_request (
	challenge varchar(40) NOT NULL PRIMARY KEY,
	verifier varchar(40) NOT NULL,
	client_id varchar(255) NOT NULL,
	subject varchar(255) NOT NULL,
	request_url text NOT NULL,
	skip bool NOT NULL,
	requested_scope text NOT NULL,
	csrf varchar(40) NOT NULL,
	authenticated_at timestamp NULL,
	requested_at timestamp NOT NULL DEFAULT now(),
	oidc_context text NOT NULL,
	forced_subject_identifier VARCHAR(255) NULL DEFAULT '',
	login_session_id VARCHAR(40) NULL,
	login_challenge VARCHAR(40) NULL,
	requested_at_audience text NULL DEFAULT '',
	acr text NULL DEFAULT '',
	context TEXT NOT NULL DEFAULT '{}',
	INDEX (client_id),
	INDEX (subject),
	INDEX (login_session_id),
	INDEX (login_challenge),
	UNIQUE INDEX (verifier)
);
CREATE TABLE hydra_oauth2_authentication_request (
	challenge varchar(40) NOT NULL PRIMARY KEY,
	requested_scope text NOT NULL,
	verifier varchar(40) NOT NULL,
	csrf varchar(40) NOT NULL,
	subject varchar(255) NOT NULL,
	request_url text NOT NULL,
	skip bool NOT NULL,
	client_id varchar(255) NOT NULL,
	requested_at timestamp NOT NULL DEFAULT now(),
	authenticated_at timestamp NULL,
	oidc_context text NOT NULL,
	login_session_id VARCHAR(40) NULL DEFAULT '',
	requested_at_audience text NULL DEFAULT '',
	INDEX (client_id),
	INDEX (subject),
	INDEX (login_session_id),
	UNIQUE INDEX (verifier)
);
CREATE TABLE hydra_oauth2_authentication_session (
	id varchar(40) NOT NULL PRIMARY KEY,
	authenticated_at timestamp NOT NULL DEFAULT NOW(),
	subject varchar(255) NOT NULL,
	remember bool NOT NULL DEFAULT FALSE
);
CREATE TABLE hydra_oauth2_consent_request_handled (
	challenge varchar(40) NOT NULL PRIMARY KEY,
	granted_scope text NOT NULL,
	remember bool NOT NULL,
	remember_for int NOT NULL,
	error text NOT NULL,
	requested_at timestamp NOT NULL DEFAULT now(),
	session_access_token text NOT NULL,
	session_id_token text NOT NULL,
	authenticated_at timestamp NULL,
	was_used bool NOT NULL,
	granted_at_audience TEXT NULL DEFAULT ''
);
CREATE TABLE hydra_oauth2_authentication_request_handled (
	challenge varchar(40) NOT NULL PRIMARY KEY,
	subject varchar(255) NOT NULL,
	remember bool NOT NULL,
	remember_for int NOT NULL,
	error text NOT NULL,
	acr text NOT NULL,
	requested_at timestamp NOT NULL DEFAULT now(),
	authenticated_at timestamp NULL,
	was_used bool NOT NULL,
	forced_subject_identifier VARCHAR(255) NULL DEFAULT '',
	context TEXT NOT NULL DEFAULT '{}'
);
CREATE TABLE hydra_oauth2_obfuscated_authentication_session (
	subject varchar(255) NOT NULL,
	client_id varchar(255) NOT NULL,
	subject_obfuscated varchar(255) NOT NULL,
	PRIMARY KEY(subject, client_id),
	INDEX (client_id, subject_obfuscated)
);
CREATE TABLE hydra_oauth2_logout_request (
	challenge varchar(36) NOT NULL PRIMARY KEY,
	verifier varchar(36) NOT NULL,
	subject varchar(255) NOT NULL,
	sid varchar(36) NOT NULL,
	client_id varchar(255),
	request_url text NOT NULL,
	redir_url text NOT NULL,
	was_used bool NOT NULL default false,
	accepted bool NOT NULL default false,
	rejected bool NOT NULL default false,
	rp_initiated bool NOT NULL default false,
	INDEX (client_id),
	UNIQUE INDEX (verifier)
);

ALTER TABLE hydra_oauth2_consent_request_handled ADD CONSTRAINT hydra_oauth2_consent_request_handled_challenge_fk FOREIGN KEY (challenge) REFERENCES hydra_oauth2_consent_request(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_authentication_request_handled ADD CONSTRAINT hydra_oauth2_authentication_request_handled_challenge_fk FOREIGN KEY (challenge) REFERENCES hydra_oauth2_authentication_request(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_consent_request ADD CONSTRAINT hydra_oauth2_consent_request_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_authentication_request ADD CONSTRAINT hydra_oauth2_authentication_request_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_obfuscated_authentication_session ADD CONSTRAINT hydra_oauth2_obfuscated_authentication_session_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_authentication_request ADD CONSTRAINT hydra_oauth2_authentication_request_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES hydra_oauth2_authentication_session(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_consent_request ADD CONSTRAINT hydra_oauth2_consent_request_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES hydra_oauth2_authentication_session(id) ON DELETE SET NULL;
ALTER TABLE hydra_oauth2_consent_request ADD CONSTRAINT hydra_oauth2_consent_request_login_challenge_fk FOREIGN KEY (login_challenge) REFERENCES hydra_oauth2_authentication_request(challenge) ON DELETE SET NULL;
ALTER TABLE hydra_oauth2_logout_request ADD CONSTRAINT hydra_oauth2_logout_request_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
