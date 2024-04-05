ALTER TABLE hydra_oauth2_oidc ADD COLUMN expires_at TIMESTAMP NULL;
ALTER TABLE hydra_oauth2_access ADD COLUMN expires_at TIMESTAMP NULL;
ALTER TABLE hydra_oauth2_refresh ADD COLUMN expires_at TIMESTAMP NULL;
ALTER TABLE hydra_oauth2_code ADD COLUMN expires_at TIMESTAMP NULL;
ALTER TABLE hydra_oauth2_pkce ADD COLUMN expires_at TIMESTAMP NULL;

CREATE INDEX hydra_oauth2_oidc_expires_at_idx ON hydra_oauth2_oidc (expires_at);
CREATE INDEX hydra_oauth2_access_expires_at_idx ON hydra_oauth2_oidc (expires_at);
CREATE INDEX hydra_oauth2_refresh_expires_at_idx ON hydra_oauth2_oidc (expires_at);
CREATE INDEX hydra_oauth2_code_expires_at_idx ON hydra_oauth2_oidc (expires_at);
CREATE INDEX hydra_oauth2_pkce_expires_at_idx ON hydra_oauth2_oidc (expires_at);
