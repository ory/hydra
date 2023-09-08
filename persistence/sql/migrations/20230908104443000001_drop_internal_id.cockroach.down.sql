ALTER TABLE hydra_client ADD COLUMN pk UUID NOT NULL DEFAULT gen_random_uuid();
