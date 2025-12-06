DROP INDEX hydra_oauth2_oidc_expires_at_idx;
DROP INDEX hydra_oauth2_access_expires_at_idx;
DROP INDEX hydra_oauth2_refresh_expires_at_idx;
DROP INDEX hydra_oauth2_code_expires_at_idx;
DROP INDEX hydra_oauth2_pkce_expires_at_idx;

ALTER TABLE hydra_oauth2_oidc DROP COLUMN expires_at;
ALTER TABLE hydra_oauth2_access DROP COLUMN expires_at;
ALTER TABLE hydra_oauth2_refresh DROP COLUMN expires_at;
ALTER TABLE hydra_oauth2_code DROP COLUMN expires_at;
ALTER TABLE hydra_oauth2_pkce DROP COLUMN expires_at;
