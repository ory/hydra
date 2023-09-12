UPDATE hydra_client SET pk = gen_random_uuid() WHERE pk IS NULL;

ALTER TABLE hydra_client ALTER COLUMN pk SET NOT NULL;
