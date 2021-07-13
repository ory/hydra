ALTER TABLE hydra_oauth2_authentication_session ALTER COLUMN authenticated_at SET DEFAULT NOW();
ALTER TABLE hydra_oauth2_authentication_session ALTER COLUMN authenticated_at SET NOT NULL;
