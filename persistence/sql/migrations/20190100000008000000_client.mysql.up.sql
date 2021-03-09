UPDATE hydra_client SET allowed_cors_origins='';
ALTER TABLE hydra_client MODIFY allowed_cors_origins TEXT NOT NULL;
