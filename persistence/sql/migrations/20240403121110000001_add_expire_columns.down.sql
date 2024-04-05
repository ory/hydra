ALTER TABLE hydra_oauth2_oidc DROP COLUMN expires_at;
ALTER TABLE hydra_oauth2_access DROP COLUMN expires_at;
ALTER TABLE hydra_oauth2_refresh DROP COLUMN expires_at;
ALTER TABLE hydra_oauth2_code DROP COLUMN expires_at;
ALTER TABLE hydra_oauth2_pkce DROP COLUMN expires_at;
