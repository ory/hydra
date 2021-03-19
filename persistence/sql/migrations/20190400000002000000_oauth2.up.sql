ALTER TABLE hydra_oauth2_access ADD subject varchar(255) NOT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_refresh ADD subject varchar(255) NOT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_code ADD subject varchar(255) NOT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_oidc ADD subject varchar(255) NOT NULL DEFAULT '';

