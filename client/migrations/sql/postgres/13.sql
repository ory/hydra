-- +migrate Up
ALTER TABLE hydra_client ADD frontchannel_logout_uri TEXT NOT NULL DEFAULT '';
ALTER TABLE hydra_client ADD frontchannel_logout_session_required BOOL NOT NULL DEFAULT FALSE;
ALTER TABLE hydra_client ADD post_logout_redirect_uris TEXT NOT NULL DEFAULT '';

ALTER TABLE hydra_client ADD backchannel_logout_uri TEXT NOT NULL DEFAULT '';
ALTER TABLE hydra_client ADD backchannel_logout_session_required BOOL NOT NULL DEFAULT FALSE;

-- +migrate Down
ALTER TABLE hydra_client DROP COLUMN frontchannel_logout_uri;
ALTER TABLE hydra_client DROP COLUMN frontchannel_logout_session_required;
ALTER TABLE hydra_client DROP COLUMN post_logout_redirect_uris;
ALTER TABLE hydra_client DROP COLUMN backchannel_logout_uri;
ALTER TABLE hydra_client DROP COLUMN backchannel_logout_session_required;
