UPDATE hydra_client SET allowed_cors_origins='';
ALTER TABLE hydra_client ALTER COLUMN allowed_cors_origins SET NOT NULL;
