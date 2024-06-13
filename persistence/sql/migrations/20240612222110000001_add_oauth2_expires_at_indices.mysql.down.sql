CREATE INDEX IF NOT EXISTS hydra_oauth2_access_expires_at_idx ON hydra_oauth2_oidc (expires_at);
CREATE INDEX IF NOT EXISTS hydra_oauth2_refresh_expires_at_idx ON hydra_oauth2_oidc (expires_at);
CREATE INDEX IF NOT EXISTS hydra_oauth2_code_expires_at_idx ON hydra_oauth2_oidc (expires_at);
CREATE INDEX IF NOT EXISTS hydra_oauth2_pkce_expires_at_idx ON hydra_oauth2_oidc (expires_at);

DROP INDEX hydra_oauth2_access_expires_at_v2_idx ON hydra_oauth2_access;
DROP INDEX hydra_oauth2_refresh_expires_at_v2_idx ON hydra_oauth2_refresh;
DROP INDEX hydra_oauth2_code_expires_at_v2_idx ON hydra_oauth2_code;
DROP INDEX hydra_oauth2_pkce_expires_at_v2_idx ON hydra_oauth2_pkce;
