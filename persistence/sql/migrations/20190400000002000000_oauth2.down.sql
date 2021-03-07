ALTER TABLE hydra_oauth2_access  DROP COLUMN subject;
ALTER TABLE hydra_oauth2_refresh DROP COLUMN subject;
ALTER TABLE hydra_oauth2_code DROP COLUMN subject;
ALTER TABLE hydra_oauth2_oidc DROP COLUMN subject;
