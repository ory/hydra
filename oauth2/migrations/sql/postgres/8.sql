-- +migrate Up
DELETE FROM hydra_oauth2_access WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_access.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_refresh WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_refresh.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_code WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_code.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_oidc WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_oidc.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_pkce WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_pkce.client_id = hydra_client.id);

DELETE FROM hydra_oauth2_access WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_consent_request_handled WHERE hydra_oauth2_access.request_id = hydra_oauth2_consent_request_handled.challenge);
DELETE FROM hydra_oauth2_refresh WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_consent_request_handled WHERE hydra_oauth2_refresh.request_id = hydra_oauth2_consent_request_handled.challenge);
DELETE FROM hydra_oauth2_code WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_consent_request_handled WHERE hydra_oauth2_code.request_id = hydra_oauth2_consent_request_handled.challenge);
DELETE FROM hydra_oauth2_oidc WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_consent_request_handled WHERE hydra_oauth2_oidc.request_id = hydra_oauth2_consent_request_handled.challenge);
DELETE FROM hydra_oauth2_pkce WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_consent_request_handled WHERE hydra_oauth2_pkce.request_id = hydra_oauth2_consent_request_handled.challenge);

-- we also want to remove all columns that have more than 40 char lengths as we're resizing in the next step
DELETE FROM hydra_oauth2_access WHERE LENGTH(request_id) > 40 OR request_id = '';
DELETE FROM hydra_oauth2_refresh WHERE LENGTH(request_id) > 40 OR request_id = '';
DELETE FROM hydra_oauth2_code WHERE LENGTH(request_id) > 40 OR request_id = '';
DELETE FROM hydra_oauth2_oidc WHERE LENGTH(request_id) > 40 OR request_id = '';
DELETE FROM hydra_oauth2_pkce WHERE LENGTH(request_id) > 40 OR request_id = '';

-- request_id is a 40 varchar in the refrerenced table
ALTER TABLE hydra_oauth2_access ALTER COLUMN request_id TYPE varchar(40);
ALTER TABLE hydra_oauth2_refresh ALTER COLUMN request_id TYPE varchar(40);
ALTER TABLE hydra_oauth2_code ALTER COLUMN request_id TYPE varchar(40);
ALTER TABLE hydra_oauth2_oidc ALTER COLUMN request_id TYPE varchar(40);
ALTER TABLE hydra_oauth2_pkce ALTER COLUMN request_id TYPE varchar(40);

-- we also want to remove all columns that have a client id with more then 255 chars
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

ALTER TABLE hydra_oauth2_access ALTER COLUMN client_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_refresh ALTER COLUMN client_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_code ALTER COLUMN client_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_oidc ALTER COLUMN client_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_pkce ALTER COLUMN client_id TYPE varchar(255);

CREATE INDEX hydra_oauth2_access_client_id_idx ON hydra_oauth2_access (client_id);
CREATE INDEX hydra_oauth2_refresh_client_id_idx ON hydra_oauth2_refresh (client_id);
CREATE INDEX hydra_oauth2_code_client_id_idx ON hydra_oauth2_code (client_id);
CREATE INDEX hydra_oauth2_oidc_client_id_idx ON hydra_oauth2_oidc (client_id);
CREATE INDEX hydra_oauth2_pkce_client_id_idx ON hydra_oauth2_pkce (client_id);

SET session_replication_role = replica;

ALTER TABLE hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;

ALTER TABLE hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_request_id_fk FOREIGN KEY (request_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_request_id_fk FOREIGN KEY (request_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_request_id_fk FOREIGN KEY (request_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_request_id_fk FOREIGN KEY (request_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_request_id_fk FOREIGN KEY (request_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;

SET session_replication_role = DEFAULT;

-- +migrate Down
ALTER TABLE hydra_oauth2_access DROP CONSTRAINT hydra_oauth2_access_client_id_fk;
ALTER TABLE hydra_oauth2_refresh DROP CONSTRAINT hydra_oauth2_refresh_client_id_fk;
ALTER TABLE hydra_oauth2_code DROP CONSTRAINT hydra_oauth2_code_client_id_fk;
ALTER TABLE hydra_oauth2_oidc DROP CONSTRAINT hydra_oauth2_oidc_client_id_fk;
ALTER TABLE hydra_oauth2_pkce DROP CONSTRAINT hydra_oauth2_pkce_client_id_fk;

ALTER TABLE hydra_oauth2_access DROP CONSTRAINT hydra_oauth2_access_request_id_fk;
ALTER TABLE hydra_oauth2_refresh DROP CONSTRAINT hydra_oauth2_refresh_request_id_fk;
ALTER TABLE hydra_oauth2_code DROP CONSTRAINT hydra_oauth2_code_request_id_fk;
ALTER TABLE hydra_oauth2_oidc DROP CONSTRAINT hydra_oauth2_oidc_request_id_fk;
ALTER TABLE hydra_oauth2_pkce DROP CONSTRAINT hydra_oauth2_pkce_request_id_fk;

ALTER TABLE hydra_oauth2_access ALTER COLUMN request_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_refresh ALTER COLUMN request_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_code ALTER COLUMN request_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_oidc ALTER COLUMN request_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_pkce ALTER COLUMN request_id TYPE varchar(255);

ALTER TABLE hydra_oauth2_access ALTER COLUMN client_id TYPE TEXT;
ALTER TABLE hydra_oauth2_refresh ALTER COLUMN client_id TYPE TEXT;
ALTER TABLE hydra_oauth2_code ALTER COLUMN client_id TYPE TEXT;
ALTER TABLE hydra_oauth2_oidc ALTER COLUMN client_id TYPE TEXT;
ALTER TABLE hydra_oauth2_pkce ALTER COLUMN client_id TYPE TEXT;
