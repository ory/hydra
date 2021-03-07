UPDATE hydra_client SET sector_identifier_uri='', jwks='', jwks_uri='', request_uris='';
ALTER TABLE hydra_client MODIFY sector_identifier_uri TEXT NOT NULL;
ALTER TABLE hydra_client MODIFY jwks TEXT NOT NULL;
ALTER TABLE hydra_client MODIFY jwks_uri TEXT NOT NULL;
ALTER TABLE hydra_client MODIFY request_uris TEXT NOT NULL;
