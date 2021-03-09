
-- Drops the client_id foreign key
ALTER TABLE hydra_oauth2_access DROP FOREIGN KEY hydra_oauth2_access_client_id_fk;
ALTER TABLE hydra_oauth2_refresh DROP FOREIGN KEY hydra_oauth2_refresh_client_id_fk;
ALTER TABLE hydra_oauth2_code DROP FOREIGN KEY hydra_oauth2_code_client_id_fk;
ALTER TABLE hydra_oauth2_oidc DROP FOREIGN KEY hydra_oauth2_oidc_client_id_fk;
ALTER TABLE hydra_oauth2_pkce DROP FOREIGN KEY hydra_oauth2_pkce_client_id_fk;

-- Drops the challenge/request_id foreign key
ALTER TABLE hydra_oauth2_access DROP FOREIGN KEY hydra_oauth2_access_challenge_id_fk;
ALTER TABLE hydra_oauth2_refresh DROP FOREIGN KEY hydra_oauth2_refresh_challenge_id_fk;
ALTER TABLE hydra_oauth2_code DROP FOREIGN KEY hydra_oauth2_code_challenge_id_fk;
ALTER TABLE hydra_oauth2_oidc DROP FOREIGN KEY hydra_oauth2_oidc_challenge_id_fk;
ALTER TABLE hydra_oauth2_pkce DROP FOREIGN KEY hydra_oauth2_pkce_challenge_id_fk;

-- Drops the index for client_id
DROP INDEX hydra_oauth2_access_client_id_idx ON hydra_oauth2_access;
DROP INDEX hydra_oauth2_refresh_client_id_idx ON hydra_oauth2_refresh;
DROP INDEX hydra_oauth2_code_client_id_idx ON hydra_oauth2_code;
DROP INDEX hydra_oauth2_oidc_client_id_idx ON hydra_oauth2_oidc;
DROP INDEX hydra_oauth2_pkce_client_id_idx ON hydra_oauth2_pkce;

-- request_id was set to varchar(40) NULL (without default) - let's revert that
ALTER TABLE hydra_oauth2_access MODIFY request_id varchar(255) NOT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_refresh MODIFY request_id varchar(255) NOT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_code MODIFY request_id varchar(255) NOT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_oidc MODIFY request_id varchar(255) NOT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_pkce MODIFY request_id varchar(255) NOT NULL DEFAULT '';

ALTER TABLE hydra_oauth2_access MODIFY client_id TEXT NOT NULL;
ALTER TABLE hydra_oauth2_refresh MODIFY client_id TEXT NOT NULL;
ALTER TABLE hydra_oauth2_code MODIFY client_id TEXT NOT NULL;
ALTER TABLE hydra_oauth2_oidc MODIFY client_id TEXT NOT NULL;
ALTER TABLE hydra_oauth2_pkce MODIFY client_id TEXT NOT NULL;