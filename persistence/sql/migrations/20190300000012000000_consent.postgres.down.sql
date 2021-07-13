DELETE FROM hydra_oauth2_logout_request WHERE client_id IS NULL;
ALTER TABLE hydra_oauth2_logout_request ALTER COLUMN client_id SET NOT NULL;
