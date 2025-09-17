ALTER TABLE hydra_oauth2_authentication_session ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP NULL;
