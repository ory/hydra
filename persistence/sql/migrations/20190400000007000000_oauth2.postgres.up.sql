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

