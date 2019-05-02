-- +migrate Up
ALTER TABLE hydra_oauth2_logout_request ALTER COLUMN client_id DROP NOT NULL;

-- +migrate Down
DELETE FROM hydra_oauth2_logout_request WHERE client_id IS NULL;
ALTER TABLE hydra_oauth2_logout_request ALTER COLUMN client_id SET NOT NULL;
