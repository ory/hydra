ALTER TABLE hydra_client ADD metadata TEXT NULL;

UPDATE hydra_client SET metadata='{}';

ALTER TABLE hydra_client ALTER COLUMN metadata SET NOT NULL;
