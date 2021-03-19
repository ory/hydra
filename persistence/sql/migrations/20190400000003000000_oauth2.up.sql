CREATE TABLE IF NOT EXISTS hydra_oauth2_pkce (
	signature      	varchar(255) NOT NULL PRIMARY KEY,
	request_id  	  varchar(255) NOT NULL,
	requested_at  	timestamp NOT NULL DEFAULT now(),
	client_id  		  text NOT NULL,
	scope  			    text NOT NULL,
	granted_scope 	text NOT NULL,
	form_data  		  text NOT NULL,
	session_data  	text NOT NULL,
	subject 		    varchar(255) NOT NULL
);

