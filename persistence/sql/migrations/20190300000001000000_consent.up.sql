CREATE TABLE hydra_oauth2_consent_request (
	challenge  			  varchar(40) NOT NULL PRIMARY KEY,
	verifier          varchar(40) NOT NULL,
	client_id			    varchar(255) NOT NULL,
	subject				    varchar(255) NOT NULL,
	request_url			  text NOT NULL,
	skip				      bool NOT NULL,
	requested_scope		text NOT NULL,
	csrf				      varchar(40) NOT NULL,
	authenticated_at	timestamp NULL,
	requested_at  		timestamp NOT NULL DEFAULT now(),
	oidc_context		  text NOT NULL
);
CREATE TABLE hydra_oauth2_authentication_request (
	challenge  			  varchar(40) NOT NULL PRIMARY KEY,
	requested_scope		text NOT NULL,
	verifier 			    varchar(40) NOT NULL,
	csrf				      varchar(40) NOT NULL,
	subject				    varchar(255) NOT NULL,
	request_url		  	text NOT NULL,
	skip				      bool NOT NULL,
	client_id			    varchar(255) NOT NULL,
	requested_at  		timestamp NOT NULL DEFAULT now(),
	authenticated_at	timestamp NULL,
	oidc_context		  text NOT NULL
);
CREATE TABLE hydra_oauth2_authentication_session (
	id      			    varchar(40) NOT NULL PRIMARY KEY,
	authenticated_at  timestamp NOT NULL DEFAULT NOW(),
	subject 			    varchar(255) NOT NULL
);
CREATE TABLE hydra_oauth2_consent_request_handled (
	challenge  				    varchar(40) NOT NULL PRIMARY KEY,
	granted_scope			    text NOT NULL,
	remember				      bool NOT NULL,
	remember_for			    int NOT NULL,
	error					        text NOT NULL,
	requested_at  			  timestamp NOT NULL DEFAULT now(),
	session_access_token 	text NOT NULL,
	session_id_token 		  text NOT NULL,
	authenticated_at		  timestamp NULL,
	was_used 				      bool NOT NULL
);
CREATE TABLE hydra_oauth2_authentication_request_handled (
	challenge  			    varchar(40) NOT NULL PRIMARY KEY,
	subject 			      varchar(255) NOT NULL,
	remember			      bool NOT NULL,
	remember_for		    int NOT NULL,
	error				        text NOT NULL,
	acr					        text NOT NULL,
	requested_at  		  timestamp NOT NULL DEFAULT now(),
	authenticated_at	  timestamp NULL,
	was_used 			      bool NOT NULL
);

