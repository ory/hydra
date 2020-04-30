-- +migrate Up
UPDATE hydra_client SET audience='';
ALTER TABLE hydra_client MODIFY audience TEXT NOT NULL;

-- +migrate Down
ALTER TABLE hydra_client MODIFY audience TEXT;
