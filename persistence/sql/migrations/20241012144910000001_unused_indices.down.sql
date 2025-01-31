-- CREATE INDEX IF NOT EXISTS hydra_oauth2_access_client_id_subject_idx ON hydra_oauth2_access (client_id ASC, subject ASC, nid ASC);
CREATE INDEX IF NOT EXISTS hydra_oauth2_access_expires_at_v2_idx ON hydra_oauth2_access (expires_at ASC);

CREATE INDEX IF NOT EXISTS hydra_oauth2_refresh_client_id_subject_idx ON hydra_oauth2_refresh (client_id ASC, subject ASC);
CREATE INDEX IF NOT EXISTS hydra_oauth2_refresh_expires_at_v2_idx ON hydra_oauth2_refresh (expires_at ASC);

CREATE INDEX IF NOT EXISTS hydra_oauth2_pkce_request_id_idx ON hydra_oauth2_pkce (request_id ASC, nid ASC);
CREATE INDEX IF NOT EXISTS hydra_oauth2_pkce_expires_at_v2_idx ON hydra_oauth2_pkce (expires_at ASC);

CREATE INDEX IF NOT EXISTS hydra_oauth2_oidc_request_id_idx ON hydra_oauth2_oidc (request_id ASC, nid ASC);
CREATE INDEX IF NOT EXISTS hydra_oauth2_oidc_expires_at_idx ON hydra_oauth2_oidc (expires_at ASC);

CREATE INDEX IF NOT EXISTS hydra_oauth2_code_request_id_idx ON hydra_oauth2_code (request_id ASC, nid ASC);
CREATE INDEX IF NOT EXISTS hydra_oauth2_code_expires_at_v2_idx ON hydra_oauth2_code (expires_at ASC);
