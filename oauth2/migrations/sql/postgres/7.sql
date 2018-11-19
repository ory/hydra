-- +migrate Up
ALTER TABLE hydra_oauth2_access ADD requested_audience TEXT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_access ADD granted_audience TEXT NULL DEFAULT '';

ALTER TABLE hydra_oauth2_refresh ADD requested_audience TEXT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_refresh ADD granted_audience TEXT NULL DEFAULT '';

ALTER TABLE hydra_oauth2_code ADD requested_audience TEXT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_code ADD granted_audience TEXT NULL DEFAULT '';

ALTER TABLE hydra_oauth2_oidc ADD requested_audience TEXT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_oidc ADD granted_audience TEXT NULL DEFAULT '';

ALTER TABLE hydra_oauth2_pkce ADD requested_audience TEXT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_pkce ADD granted_audience TEXT NULL DEFAULT '';

-- +migrate Down
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
