ALTER TABLE hydra_client ALTER PRIMARY KEY USING COLUMNS (id, nid) USING HASH WITH bucket_count=16;
