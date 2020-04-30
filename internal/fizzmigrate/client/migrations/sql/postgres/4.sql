-- +migrate Up
UPDATE hydra_client SET sector_identifier_uri='', jwks='', jwks_uri='', request_uris='';
ALTER TABLE hydra_client ALTER COLUMN sector_identifier_uri SET NOT NULL;
ALTER TABLE hydra_client ALTER COLUMN jwks SET NOT NULL;
ALTER TABLE hydra_client ALTER COLUMN jwks_uri SET NOT NULL;
ALTER TABLE hydra_client ALTER COLUMN request_uris SET NOT NULL;

-- +migrate Down
ALTER TABLE hydra_client ALTER COLUMN sector_identifier_uri DROP NOT NULL;
ALTER TABLE hydra_client ALTER COLUMN jwks DROP NOT NULL;
ALTER TABLE hydra_client ALTER COLUMN jwks_uri DROP NOT NULL;
ALTER TABLE hydra_client ALTER COLUMN request_uris DROP NOT NULL;
