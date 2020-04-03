UPDATE hydra_client SET token_endpoint_auth_method='none' WHERE public=TRUE;
ALTER TABLE hydra_client DROP COLUMN public;
