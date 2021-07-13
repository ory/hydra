ALTER TABLE hydra_oauth2_consent_request ADD forced_subject_identifier VARCHAR(255) NULL DEFAULT '';
ALTER TABLE hydra_oauth2_authentication_request_handled ADD forced_subject_identifier VARCHAR(255) NULL DEFAULT '';
CREATE TABLE hydra_oauth2_obfuscated_authentication_session (
	subject  			        varchar(255) NOT NULL,
	client_id 			      varchar(255) NOT NULL,
	subject_obfuscated	  varchar(255) NOT NULL,
	PRIMARY KEY(subject, client_id)
);

