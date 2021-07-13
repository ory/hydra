CREATE TABLE IF NOT EXISTS hydra_client (
	id      	      varchar(255) NOT NULL PRIMARY KEY,
	client_name  	  text NOT NULL,
	client_secret  	text NOT NULL,
	redirect_uris  	text NOT NULL,
	grant_types  	text NOT NULL,
	response_types  text NOT NULL,
	scope  			    text NOT NULL,
	owner  			    text NOT NULL,
	policy_uri  	  text NOT NULL,
	tos_uri  		    text NOT NULL,
	client_uri  	  text NOT NULL,
	logo_uri  		  text NOT NULL,
	contacts  		  text NOT NULL,
	public  		    boolean NOT NULL
);
