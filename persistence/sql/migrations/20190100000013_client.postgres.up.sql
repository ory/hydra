ALTER TABLE hydra_client ADD frontchannel_logout_uri TEXT NOT NULL DEFAULT '';
ALTER TABLE hydra_client ADD frontchannel_logout_session_required BOOL NOT NULL DEFAULT FALSE;
ALTER TABLE hydra_client ADD post_logout_redirect_uris TEXT NOT NULL DEFAULT '';

ALTER TABLE hydra_client ADD backchannel_logout_uri TEXT NOT NULL DEFAULT '';
ALTER TABLE hydra_client ADD backchannel_logout_session_required BOOL NOT NULL DEFAULT FALSE;
