-- hydra_client
ALTER TABLE hydra_client ADD COLUMN "nid" UUID;
ALTER TABLE hydra_client ADD CONSTRAINT "hydra_client_nid_fk_idx" FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_client SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_client ALTER nid SET NOT NULL;
--split
DROP INDEX hydra_client_id_key CASCADE;
--split
CREATE UNIQUE INDEX hydra_client_id_key ON hydra_client (id ASC, nid ASC);
--split

-- hydra_oauth2_access
ALTER TABLE hydra_oauth2_access ADD COLUMN "nid" UUID;
ALTER TABLE hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_nid_fk_idx FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_access SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_oauth2_access ALTER nid SET NOT NULL;
--split
ALTER TABLE hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES hydra_client(id, nid) ON DELETE CASCADE;
--split
DROP INDEX hydra_oauth2_access_requested_at_idx;
DROP INDEX hydra_oauth2_access_client_id_idx;
DROP INDEX hydra_oauth2_access_challenge_id_idx;
DROP INDEX hydra_oauth2_access_client_id_subject_idx;
DROP INDEX hydra_oauth2_access_request_id_idx;
--split
CREATE INDEX hydra_oauth2_access_requested_at_idx ON hydra_oauth2_access (requested_at, nid);
CREATE INDEX hydra_oauth2_access_client_id_idx ON hydra_oauth2_access (client_id, nid);
CREATE INDEX hydra_oauth2_access_challenge_id_idx ON hydra_oauth2_access (challenge_id);
CREATE INDEX hydra_oauth2_access_client_id_subject_idx ON hydra_oauth2_access (client_id, subject, nid);
CREATE INDEX hydra_oauth2_access_request_id_idx ON hydra_oauth2_access (request_id, nid);
--split

-- hydra_oauth2_authentication_session
ALTER TABLE hydra_oauth2_authentication_session ADD COLUMN "nid" UUID;
ALTER TABLE hydra_oauth2_authentication_session ADD CONSTRAINT hydra_oauth2_authentication_session_nid_fk_idx FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_authentication_session SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_oauth2_authentication_session ALTER nid SET NOT NULL;
--split
DROP INDEX hydra_oauth2_authentication_session_subject_idx;
--split
CREATE INDEX hydra_oauth2_authentication_session_subject_idx ON hydra_oauth2_authentication_session (subject ASC, nid ASC);
--split

-- hydra_oauth2_code
ALTER TABLE hydra_oauth2_code ADD COLUMN "nid" UUID;
ALTER TABLE hydra_oauth2_code ADD CONSTRAINT "hydra_oauth2_code_nid_fk_idx" FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_code SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_oauth2_code ALTER nid SET NOT NULL;
--split
ALTER TABLE hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES hydra_client(id, nid) ON DELETE CASCADE;
--split
DROP INDEX hydra_oauth2_code_client_id_idx;
DROP INDEX hydra_oauth2_code_challenge_id_idx;
DROP INDEX hydra_oauth2_code_request_id_idx;
--split
CREATE INDEX hydra_oauth2_code_client_id_idx ON hydra_oauth2_code (client_id, nid);
CREATE INDEX hydra_oauth2_code_challenge_id_idx ON hydra_oauth2_code (challenge_id, nid);
CREATE INDEX hydra_oauth2_code_request_id_idx ON hydra_oauth2_code (request_id, nid);
--split

-- hydra_oauth2_flow
ALTER TABLE hydra_oauth2_flow ADD COLUMN "nid" UUID;
ALTER TABLE hydra_oauth2_flow ADD CONSTRAINT "hydra_oauth2_flow_nid_fk_idx" FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_flow SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_oauth2_flow ALTER nid SET NOT NULL;
--split
ALTER TABLE hydra_oauth2_flow ADD CONSTRAINT hydra_oauth2_flow_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES hydra_client(id, nid) ON DELETE CASCADE;
--split
DROP INDEX hydra_oauth2_flow_client_id_subject_idx;
DROP INDEX hydra_oauth2_flow_cid_idx;
DROP INDEX hydra_oauth2_flow_login_session_id_idx;
DROP INDEX hydra_oauth2_flow_sub_idx;
DROP INDEX hydra_oauth2_flow_login_verifier_idx;
DROP INDEX hydra_oauth2_flow_consent_verifier_idx;
--split
CREATE INDEX hydra_oauth2_flow_client_id_subject_idx ON hydra_oauth2_flow (client_id ASC, nid ASC, subject ASC);
CREATE INDEX hydra_oauth2_flow_cid_idx ON hydra_oauth2_flow (client_id ASC, nid ASC);
CREATE INDEX hydra_oauth2_flow_login_session_id_idx ON hydra_oauth2_flow (login_session_id ASC, nid ASC);
CREATE INDEX hydra_oauth2_flow_sub_idx ON hydra_oauth2_flow (subject ASC, nid ASC);
CREATE UNIQUE INDEX hydra_oauth2_flow_login_verifier_idx ON hydra_oauth2_flow (login_verifier ASC);
CREATE UNIQUE INDEX hydra_oauth2_flow_consent_verifier_idx ON hydra_oauth2_flow (consent_verifier ASC);
--split

-- hydra_oauth2_jti_blacklist
ALTER TABLE hydra_oauth2_jti_blacklist ADD COLUMN "nid" UUID;
ALTER TABLE hydra_oauth2_jti_blacklist ADD CONSTRAINT "hydra_oauth2_jti_blacklist_nid_fk_idx" FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_jti_blacklist SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_oauth2_jti_blacklist ALTER nid SET NOT NULL;
--split
DROP INDEX hydra_oauth2_jti_blacklist_expires_at_idx;
--split
CREATE INDEX hydra_oauth2_jti_blacklist_expires_at_idx ON hydra_oauth2_jti_blacklist (expires_at ASC, nid ASC);
--split
ALTER TABLE hydra_oauth2_jti_blacklist
    DROP CONSTRAINT "primary",
    ADD CONSTRAINT hydra_oauth2_jti_blacklist_pkey PRIMARY KEY (signature ASC, nid ASC);
--split

-- hydra_oauth2_logout_request
ALTER TABLE hydra_oauth2_logout_request ADD COLUMN "nid" UUID;
ALTER TABLE hydra_oauth2_logout_request ADD CONSTRAINT hydra_oauth2_logout_request_nid_fk_idx FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_logout_request SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_oauth2_logout_request ALTER nid SET NOT NULL;
--split
ALTER TABLE hydra_oauth2_logout_request ADD CONSTRAINT hydra_oauth2_logout_request_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES hydra_client(id, nid) ON DELETE CASCADE;
--split
DROP INDEX hydra_oauth2_logout_request_client_id_idx;
--split
CREATE INDEX hydra_oauth2_logout_request_client_id_idx ON hydra_oauth2_logout_request (client_id ASC, nid ASC);
--split

-- hydra_oauth2_obfuscated_authentication_session
ALTER TABLE hydra_oauth2_obfuscated_authentication_session ADD COLUMN "nid" UUID;
ALTER TABLE hydra_oauth2_obfuscated_authentication_session ADD CONSTRAINT hydra_oauth2_obfuscated_authentication_session_nid_fk_idx FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_obfuscated_authentication_session SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_oauth2_obfuscated_authentication_session ALTER nid SET NOT NULL;
--split
ALTER TABLE hydra_oauth2_obfuscated_authentication_session ADD CONSTRAINT hydra_oauth2_obfuscated_authentication_session_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES hydra_client(id, nid) ON DELETE CASCADE;
--split
ALTER TABLE hydra_oauth2_obfuscated_authentication_session
    DROP CONSTRAINT "primary",
    ADD CONSTRAINT "hydra_oauth2_obfuscated_authentication_session_pkey" PRIMARY KEY (subject ASC, client_id ASC, nid ASC);
--split
DROP INDEX hydra_oauth2_obfuscated_authentication_session_client_id_subject_obfuscated_idx;
--split
CREATE UNIQUE INDEX hydra_oauth2_obfuscated_authentication_session_client_id_subject_obfuscated_idx ON hydra_oauth2_obfuscated_authentication_session (client_id ASC, subject_obfuscated ASC, nid ASC);
--split

-- hydra_oauth2_oidc
ALTER TABLE hydra_oauth2_oidc ADD COLUMN "nid" UUID;
ALTER TABLE hydra_oauth2_oidc ADD CONSTRAINT "hydra_oauth2_oidc_nid_fk_idx" FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_oidc SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_oauth2_oidc ALTER nid SET NOT NULL;
--split
ALTER TABLE hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES hydra_client(id, nid) ON DELETE CASCADE;
--split
DROP INDEX hydra_oauth2_oidc_client_id_idx;
DROP INDEX hydra_oauth2_oidc_challenge_id_idx;
DROP INDEX hydra_oauth2_oidc_request_id_idx;
--split
CREATE INDEX hydra_oauth2_oidc_client_id_idx ON hydra_oauth2_oidc (client_id ASC, nid ASC);
CREATE INDEX hydra_oauth2_oidc_challenge_id_idx ON hydra_oauth2_oidc (challenge_id ASC);
CREATE INDEX hydra_oauth2_oidc_request_id_idx ON hydra_oauth2_oidc (request_id ASC, nid ASC);
--split

-- hydra_oauth2_pkce
ALTER TABLE hydra_oauth2_pkce ADD COLUMN "nid" UUID;
ALTER TABLE hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_nid_fk_idx FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_pkce SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_oauth2_pkce ALTER nid SET NOT NULL;
--split
ALTER TABLE hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES hydra_client(id, nid) ON DELETE CASCADE;
--split
DROP INDEX hydra_oauth2_pkce_client_id_idx;
DROP INDEX hydra_oauth2_pkce_challenge_id_idx;
DROP INDEX hydra_oauth2_pkce_request_id_idx;
--split
CREATE INDEX hydra_oauth2_pkce_client_id_idx ON hydra_oauth2_pkce (client_id ASC, nid ASC);
CREATE INDEX hydra_oauth2_pkce_challenge_id_idx ON hydra_oauth2_pkce (challenge_id ASC);
CREATE INDEX hydra_oauth2_pkce_request_id_idx ON hydra_oauth2_pkce (request_id ASC, nid ASC);
--split

-- hydra_oauth2_refresh
ALTER TABLE hydra_oauth2_refresh ADD COLUMN "nid" UUID;
ALTER TABLE hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_nid_fk_idx FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_refresh SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_oauth2_refresh ALTER nid SET NOT NULL;
--split
ALTER TABLE hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES hydra_client(id, nid) ON DELETE CASCADE;
--split
DROP INDEX hydra_oauth2_refresh_client_id_idx;
DROP INDEX hydra_oauth2_refresh_challenge_id_idx;
DROP INDEX hydra_oauth2_refresh_client_id_subject_idx;
DROP INDEX hydra_oauth2_refresh_request_id_idx;
--split
CREATE INDEX hydra_oauth2_refresh_client_id_idx ON hydra_oauth2_refresh (client_id ASC, nid ASC);
CREATE INDEX hydra_oauth2_refresh_challenge_id_idx ON hydra_oauth2_refresh (challenge_id ASC);
CREATE INDEX hydra_oauth2_refresh_client_id_subject_idx ON hydra_oauth2_refresh (client_id ASC, subject ASC);
CREATE INDEX hydra_oauth2_refresh_request_id_idx ON hydra_oauth2_refresh (request_id ASC);
--split

-- hydra_jwk
ALTER TABLE hydra_jwk ADD COLUMN "nid" UUID;
ALTER TABLE hydra_jwk ADD CONSTRAINT hydra_jwk_nid_fk_idx FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_jwk SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_jwk ALTER nid SET NOT NULL;
--split
DROP INDEX hydra_jwk_sid_kid_key CASCADE;
--split
CREATE UNIQUE INDEX hydra_jwk_sid_kid_nid_key ON hydra_jwk (sid ASC, kid ASC, nid ASC);
--split

-- hydra_oauth2_trusted_jwt_bearer_issuer
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer ADD COLUMN "nid" UUID;
--split
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer ADD CONSTRAINT hydra_oauth2_trusted_jwt_bearer_issuer_nid_fk_idx FOREIGN KEY ("nid") REFERENCES "networks" ("id") ON UPDATE RESTRICT ON DELETE CASCADE;
--split
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer ADD CONSTRAINT fk_key_set_ref_hydra_jwk FOREIGN KEY (key_set, key_id, nid) REFERENCES hydra_jwk(sid, kid, nid) ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_trusted_jwt_bearer_issuer SET nid = (SELECT id FROM networks LIMIT 1);
--split
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer ALTER nid SET NOT NULL;
--split
DROP INDEX hydra_oauth2_trusted_jwt_bearer_issuer_issuer_subject_key_id_key CASCADE;
DROP INDEX hydra_oauth2_trusted_jwt_bearer_issuer_expires_at_idx;
--split
CREATE UNIQUE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_issuer_subject_key_id_key ON hydra_oauth2_trusted_jwt_bearer_issuer (issuer ASC, subject ASC, key_id ASC, nid ASC);
CREATE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_expires_at_idx ON hydra_oauth2_trusted_jwt_bearer_issuer (expires_at ASC);
CREATE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_nid_idx ON hydra_oauth2_trusted_jwt_bearer_issuer (id, nid);
