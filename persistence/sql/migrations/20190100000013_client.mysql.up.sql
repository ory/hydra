ALTER TABLE hydra_client ADD frontchannel_logout_uri TEXT NULL;
ALTER TABLE hydra_client ADD frontchannel_logout_session_required BOOL NOT NULL DEFAULT FALSE;
ALTER TABLE hydra_client ADD post_logout_redirect_uris TEXT NULL;

ALTER TABLE hydra_client ADD backchannel_logout_uri TEXT NULL;
ALTER TABLE hydra_client ADD backchannel_logout_session_required BOOL NOT NULL DEFAULT FALSE;

UPDATE hydra_client SET frontchannel_logout_uri='',post_logout_redirect_uris='',backchannel_logout_uri='';

ALTER TABLE hydra_client MODIFY frontchannel_logout_uri TEXT NOT NULL;
ALTER TABLE hydra_client MODIFY post_logout_redirect_uris TEXT NOT NULL;
ALTER TABLE hydra_client MODIFY backchannel_logout_uri TEXT NOT NULL;
