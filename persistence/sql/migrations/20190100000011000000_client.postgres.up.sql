UPDATE hydra_client SET audience='';
ALTER TABLE hydra_client ALTER COLUMN audience SET NOT NULL;
