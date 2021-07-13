
-- Drops the client_id foreign key
ALTER TABLE hydra_oauth2_access DROP CONSTRAINT hydra_oauth2_access_client_id_fk;
ALTER TABLE hydra_oauth2_refresh DROP CONSTRAINT hydra_oauth2_refresh_client_id_fk;
ALTER TABLE hydra_oauth2_code DROP CONSTRAINT hydra_oauth2_code_client_id_fk;
ALTER TABLE hydra_oauth2_oidc DROP CONSTRAINT hydra_oauth2_oidc_client_id_fk;
ALTER TABLE hydra_oauth2_pkce DROP CONSTRAINT hydra_oauth2_pkce_client_id_fk;

-- Drops the challenge/request_id foreign key
ALTER TABLE hydra_oauth2_access DROP CONSTRAINT hydra_oauth2_access_challenge_id_fk;
ALTER TABLE hydra_oauth2_refresh DROP CONSTRAINT hydra_oauth2_refresh_challenge_id_fk;
ALTER TABLE hydra_oauth2_code DROP CONSTRAINT hydra_oauth2_code_challenge_id_fk;
ALTER TABLE hydra_oauth2_oidc DROP CONSTRAINT hydra_oauth2_oidc_challenge_id_fk;
ALTER TABLE hydra_oauth2_pkce DROP CONSTRAINT hydra_oauth2_pkce_challenge_id_fk;

-- Drops the index for client_id
DROP INDEX hydra_oauth2_access_client_id_idx;
DROP INDEX hydra_oauth2_refresh_client_id_idx;
DROP INDEX hydra_oauth2_code_client_id_idx;
DROP INDEX hydra_oauth2_oidc_client_id_idx;
DROP INDEX hydra_oauth2_pkce_client_id_idx;

-- request_id was set to varchar(40) NULL (without default) - let's revert that
ALTER TABLE hydra_oauth2_access ALTER COLUMN request_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_refresh ALTER COLUMN request_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_code ALTER COLUMN request_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_oidc ALTER COLUMN request_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_pkce ALTER COLUMN request_id TYPE varchar(255);

-- client_id was set to varchar(255), let's revert that.
ALTER TABLE hydra_oauth2_access ALTER COLUMN client_id TYPE TEXT;
ALTER TABLE hydra_oauth2_refresh ALTER COLUMN client_id TYPE TEXT;
ALTER TABLE hydra_oauth2_code ALTER COLUMN client_id TYPE TEXT;
ALTER TABLE hydra_oauth2_oidc ALTER COLUMN client_id TYPE TEXT;
ALTER TABLE hydra_oauth2_pkce ALTER COLUMN client_id TYPE TEXT;
