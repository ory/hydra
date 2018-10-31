-- +migrate Up
UPDATE hydra_client SET token_endpoint_auth_method='none' WHERE public=TRUE;
ALTER TABLE hydra_client DROP COLUMN public;

-- +migrate Down
ALTER TABLE hydra_client ADD public BOOLEAN NOT NULL DEFAULT FALSE;
UPDATE hydra_client SET public=TRUE WHERE token_endpoint_auth_method='none';
