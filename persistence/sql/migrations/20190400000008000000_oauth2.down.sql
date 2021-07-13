ALTER TABLE hydra_oauth2_access  DROP COLUMN challenge_id;
ALTER TABLE hydra_oauth2_refresh DROP COLUMN challenge_id;
ALTER TABLE hydra_oauth2_code DROP COLUMN challenge_id;
ALTER TABLE hydra_oauth2_oidc DROP COLUMN challenge_id;
ALTER TABLE hydra_oauth2_pkce DROP COLUMN challenge_id;
