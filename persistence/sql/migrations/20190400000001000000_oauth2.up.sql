CREATE TABLE IF NOT EXISTS hydra_oauth2_access (
	signature      	varchar(255) NOT NULL PRIMARY KEY,
	request_id  	  varchar(255) NOT NULL,
	requested_at  	timestamp NOT NULL DEFAULT now(),
	client_id  		  text NOT NULL,
	scope  			    text NOT NULL,
	granted_scope 	text NOT NULL,
	form_data  		  text NOT NULL,
	session_data  	text NOT NULL
);

CREATE TABLE IF NOT EXISTS hydra_oauth2_refresh (
	signature      	varchar(255) NOT NULL PRIMARY KEY,
	request_id  	  varchar(255) NOT NULL,
	requested_at  	timestamp NOT NULL DEFAULT now(),
	client_id  		  text NOT NULL,
	scope  			    text NOT NULL,
	granted_scope 	text NOT NULL,
	form_data  		  text NOT NULL,
	session_data  	text NOT NULL
);

CREATE TABLE IF NOT EXISTS hydra_oauth2_code (
	signature      	varchar(255) NOT NULL PRIMARY KEY,
	request_id  	  varchar(255) NOT NULL,
	requested_at  	timestamp NOT NULL DEFAULT now(),
	client_id  		  text NOT NULL,
	scope  			    text NOT NULL,
	granted_scope 	text NOT NULL,
	form_data  		  text NOT NULL,
	session_data  	text NOT NULL
);

CREATE TABLE IF NOT EXISTS hydra_oauth2_oidc (
	signature      	varchar(255) NOT NULL PRIMARY KEY,
	request_id  	  varchar(255) NOT NULL,
	requested_at  	timestamp NOT NULL DEFAULT now(),
	client_id  		  text NOT NULL,
	scope  			    text NOT NULL,
	granted_scope 	text NOT NULL,
	form_data  		  text NOT NULL,
	session_data  	text NOT NULL
);

