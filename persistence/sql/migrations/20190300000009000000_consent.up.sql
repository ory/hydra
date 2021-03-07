CREATE TABLE hydra_oauth2_logout_request (
	challenge  			        varchar(36) NOT NULL PRIMARY KEY,
	verifier  			        varchar(36) NOT NULL,
	subject				          varchar(255) NOT NULL,
	sid   				          varchar(36) NOT NULL,
	client_id               varchar(255) NOT NULL,
	request_url		  	      text NOT NULL,
	redir_url 		  	      text NOT NULL,
	was_used                bool NOT NULL default false,
  accepted                bool NOT NULL default false,
  rejected                bool NOT NULL default false,
	rp_initiated            bool NOT NULL default false
);

