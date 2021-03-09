ALTER TABLE hydra_oauth2_access DROP COLUMN requested_audience;
ALTER TABLE hydra_oauth2_access DROP COLUMN granted_audience;

ALTER TABLE hydra_oauth2_refresh DROP COLUMN requested_audience;
ALTER TABLE hydra_oauth2_refresh DROP COLUMN granted_audience;

ALTER TABLE hydra_oauth2_code DROP COLUMN requested_audience;
ALTER TABLE hydra_oauth2_code DROP COLUMN granted_audience;

ALTER TABLE hydra_oauth2_oidc DROP COLUMN requested_audience;
ALTER TABLE hydra_oauth2_oidc DROP COLUMN granted_audience;

ALTER TABLE hydra_oauth2_pkce DROP COLUMN requested_audience;
ALTER TABLE hydra_oauth2_pkce DROP COLUMN granted_audience;
