UPDATE hydra_client SET pk = gen_random_uuid() WHERE pk IS NULL;

ALTER TABLE hydra_client ALTER COLUMN pk SET NOT NULL;

ALTER TABLE hydra_client DROP CONSTRAINT hydra_client_pkey;

ALTER TABLE hydra_client ADD PRIMARY KEY (pk);
