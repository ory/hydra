ALTER TABLE hydra_client ALTER COLUMN sector_identifier_uri DROP NOT NULL;
ALTER TABLE hydra_client ALTER COLUMN jwks DROP NOT NULL;
ALTER TABLE hydra_client ALTER COLUMN jwks_uri DROP NOT NULL;
ALTER TABLE hydra_client ALTER COLUMN request_uris DROP NOT NULL;
