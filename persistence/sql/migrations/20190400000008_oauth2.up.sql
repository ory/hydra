ALTER TABLE hydra_oauth2_access ADD challenge_id varchar(40) NULL;
ALTER TABLE hydra_oauth2_refresh ADD challenge_id varchar(40) NULL;
ALTER TABLE hydra_oauth2_code ADD challenge_id varchar(40) NULL;
ALTER TABLE hydra_oauth2_oidc ADD challenge_id varchar(40) NULL;
ALTER TABLE hydra_oauth2_pkce ADD challenge_id varchar(40) NULL;

