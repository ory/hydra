ALTER TABLE hydra_oauth2_authentication_session ALTER COLUMN authenticated_at DROP NOT NULL;
ALTER TABLE hydra_oauth2_authentication_session ALTER COLUMN authenticated_at DROP DEFAULT;
