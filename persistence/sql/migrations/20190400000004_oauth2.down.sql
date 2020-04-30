ALTER TABLE hydra_oauth2_access DROP COLUMN active;
ALTER TABLE hydra_oauth2_refresh DROP COLUMN active;
ALTER TABLE hydra_oauth2_code DROP COLUMN active;
ALTER TABLE hydra_oauth2_oidc DROP COLUMN active;
ALTER TABLE hydra_oauth2_pkce DROP COLUMN active;
