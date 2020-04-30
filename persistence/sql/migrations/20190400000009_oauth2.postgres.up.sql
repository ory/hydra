
-- First we need to delete all rows that point to a non-existing oauth2 client.
DELETE FROM hydra_oauth2_access WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_access.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_refresh WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_refresh.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_code WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_code.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_oidc WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_oidc.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_pkce WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_pkce.client_id = hydra_client.id);

-- request_id is a 40 varchar in the referenced table which is why we are resizing
-- 1. We must remove request_ids longer than 40 chars. This should never happen as we've never issued them longer than this
DELETE FROM hydra_oauth2_access WHERE LENGTH(request_id) > 40;
DELETE FROM hydra_oauth2_refresh WHERE LENGTH(request_id) > 40;
DELETE FROM hydra_oauth2_code WHERE LENGTH(request_id) > 40;
DELETE FROM hydra_oauth2_oidc WHERE LENGTH(request_id) > 40;
DELETE FROM hydra_oauth2_pkce WHERE LENGTH(request_id) > 40;

-- 2. Next we're actually resizing
ALTER TABLE hydra_oauth2_access ALTER COLUMN request_id TYPE varchar(40);
ALTER TABLE hydra_oauth2_refresh ALTER COLUMN request_id TYPE varchar(40);
ALTER TABLE hydra_oauth2_code ALTER COLUMN request_id TYPE varchar(40);
ALTER TABLE hydra_oauth2_oidc ALTER COLUMN request_id TYPE varchar(40);
ALTER TABLE hydra_oauth2_pkce ALTER COLUMN request_id TYPE varchar(40);

-- In preparation for creating the client_id index and foreign key, we must set it to varchar(255) which is also
-- the length of hydra_client.id
DELETE FROM hydra_oauth2_access WHERE LENGTH(client_id) > 255;
DELETE FROM hydra_oauth2_refresh WHERE LENGTH(client_id) > 255;
DELETE FROM hydra_oauth2_code WHERE LENGTH(client_id) > 255;
DELETE FROM hydra_oauth2_oidc WHERE LENGTH(client_id) > 255;
DELETE FROM hydra_oauth2_pkce WHERE LENGTH(client_id) > 255;
ALTER TABLE hydra_oauth2_access ALTER COLUMN client_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_refresh ALTER COLUMN client_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_code ALTER COLUMN client_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_oidc ALTER COLUMN client_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_pkce ALTER COLUMN client_id TYPE varchar(255);

-- Now it's time to create the index for client_id
CREATE INDEX hydra_oauth2_access_client_id_idx ON hydra_oauth2_access (client_id);
CREATE INDEX hydra_oauth2_refresh_client_id_idx ON hydra_oauth2_refresh (client_id);
CREATE INDEX hydra_oauth2_code_client_id_idx ON hydra_oauth2_code (client_id);
CREATE INDEX hydra_oauth2_oidc_client_id_idx ON hydra_oauth2_oidc (client_id);
CREATE INDEX hydra_oauth2_pkce_client_id_idx ON hydra_oauth2_pkce (client_id);

-- Foreign keys start here

-- This caused #1209:
-- SET session_replication_role = replica;

-- This creates a foreign key that cascade delete's if the client_id is removed.
ALTER TABLE hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;

-- This creates a foreign key that cascade delete's if the consent associated with this is removed.
ALTER TABLE hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;

-- This caused #1209:
-- SET session_replication_role = DEFAULT;

