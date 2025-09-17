ALTER TABLE hydra_oauth2_refresh ADD COLUMN IF NOT EXISTS access_token_signature VARCHAR(255) DEFAULT NULL;
