ALTER TABLE hydra_client DROP CONSTRAINT hydra_client_pkey;

ALTER TABLE hydra_client ALTER COLUMN pk DROP NOT NULL;

ALTER TABLE hydra_client ADD PRIMARY KEY (id, nid);
