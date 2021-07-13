CREATE INDEX hydra_oauth2_access_challenge_id_idx ON hydra_oauth2_access (challenge_id);
CREATE INDEX hydra_oauth2_refresh_challenge_id_idx ON hydra_oauth2_refresh (challenge_id);
CREATE INDEX hydra_oauth2_code_challenge_id_idx ON hydra_oauth2_code (challenge_id);
CREATE INDEX hydra_oauth2_oidc_challenge_id_idx ON hydra_oauth2_oidc (challenge_id);
CREATE INDEX hydra_oauth2_pkce_challenge_id_idx ON hydra_oauth2_pkce (challenge_id);

