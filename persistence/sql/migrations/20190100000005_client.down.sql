ALTER TABLE hydra_client ADD public BOOLEAN NOT NULL DEFAULT FALSE;
UPDATE hydra_client SET public=TRUE WHERE token_endpoint_auth_method='none';
