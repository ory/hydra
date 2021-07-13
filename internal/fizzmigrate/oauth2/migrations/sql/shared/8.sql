-- +migrate Up
ALTER TABLE hydra_oauth2_access ADD challenge_id varchar(40) NULL;
ALTER TABLE hydra_oauth2_refresh ADD challenge_id varchar(40) NULL;
ALTER TABLE hydra_oauth2_code ADD challenge_id varchar(40) NULL;
ALTER TABLE hydra_oauth2_oidc ADD challenge_id varchar(40) NULL;
ALTER TABLE hydra_oauth2_pkce ADD challenge_id varchar(40) NULL;

-- +migrate Down
ALTER TABLE hydra_oauth2_access  DROP COLUMN challenge_id;
ALTER TABLE hydra_oauth2_refresh DROP COLUMN challenge_id;
ALTER TABLE hydra_oauth2_code DROP COLUMN challenge_id;
ALTER TABLE hydra_oauth2_oidc DROP COLUMN challenge_id;
ALTER TABLE hydra_oauth2_pkce DROP COLUMN challenge_id;
