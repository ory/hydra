DROP INDEX hydra_oauth2_access_expires_at_idx;
DROP INDEX hydra_oauth2_refresh_expires_at_idx;
DROP INDEX hydra_oauth2_code_expires_at_idx;
DROP INDEX hydra_oauth2_pkce_expires_at_idx;

CREATE INDEX hydra_oauth2_access_expires_at_idx ON hydra_oauth2_access (expires_at);
CREATE INDEX hydra_oauth2_refresh_expires_at_idx ON hydra_oauth2_refresh (expires_at);
CREATE INDEX hydra_oauth2_code_expires_at_idx ON hydra_oauth2_code (expires_at);
CREATE INDEX hydra_oauth2_pkce_expires_at_idx ON hydra_oauth2_pkce (expires_at);
