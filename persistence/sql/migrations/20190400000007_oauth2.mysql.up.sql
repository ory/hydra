ALTER TABLE hydra_oauth2_access ADD requested_audience TEXT NULL;
UPDATE hydra_oauth2_access SET requested_audience='';
ALTER TABLE hydra_oauth2_access MODIFY requested_audience TEXT NOT NULL;
ALTER TABLE hydra_oauth2_access ADD granted_audience TEXT NULL;
UPDATE hydra_oauth2_access SET granted_audience='';
ALTER TABLE hydra_oauth2_access MODIFY granted_audience TEXT NOT NULL;

ALTER TABLE hydra_oauth2_refresh ADD requested_audience TEXT NULL;
UPDATE hydra_oauth2_refresh SET requested_audience='';
ALTER TABLE hydra_oauth2_refresh MODIFY requested_audience TEXT NOT NULL;
ALTER TABLE hydra_oauth2_refresh ADD granted_audience TEXT NULL;
UPDATE hydra_oauth2_refresh SET granted_audience='';
ALTER TABLE hydra_oauth2_refresh MODIFY granted_audience TEXT NOT NULL;

ALTER TABLE hydra_oauth2_code ADD requested_audience TEXT NULL;
UPDATE hydra_oauth2_code SET requested_audience='';
ALTER TABLE hydra_oauth2_code MODIFY requested_audience TEXT NOT NULL;
ALTER TABLE hydra_oauth2_code ADD granted_audience TEXT NULL;
UPDATE hydra_oauth2_code SET granted_audience='';
ALTER TABLE hydra_oauth2_code MODIFY granted_audience TEXT NOT NULL;

ALTER TABLE hydra_oauth2_oidc ADD requested_audience TEXT NULL;
UPDATE hydra_oauth2_oidc SET requested_audience='';
ALTER TABLE hydra_oauth2_oidc MODIFY requested_audience TEXT NOT NULL;
ALTER TABLE hydra_oauth2_oidc ADD granted_audience TEXT NULL;
UPDATE hydra_oauth2_oidc SET granted_audience='';
ALTER TABLE hydra_oauth2_oidc MODIFY granted_audience TEXT NOT NULL;

ALTER TABLE hydra_oauth2_pkce ADD requested_audience TEXT NULL;
UPDATE hydra_oauth2_pkce SET requested_audience='';
ALTER TABLE hydra_oauth2_pkce MODIFY requested_audience TEXT NOT NULL;
ALTER TABLE hydra_oauth2_pkce ADD granted_audience TEXT NULL;
UPDATE hydra_oauth2_pkce SET granted_audience='';
ALTER TABLE hydra_oauth2_pkce MODIFY granted_audience TEXT NOT NULL;

