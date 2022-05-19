-- Encode key_id in ascii as a workaround for the 3072-byte index entry size limit[1]
-- This is a breaking change for MySQL key IDs with utf-8 symbols higher than 127
-- [1]: https://dev.mysql.com/doc/refman/8.0/en/innodb-limits.html
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer DROP FOREIGN KEY `hydra_oauth2_trusted_jwt_bearer_issuer_ibfk_1`;
ALTER TABLE hydra_jwk MODIFY `kid` varchar(255) CHARACTER SET 'ascii' NOT NULL;
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer MODIFY `key_id` varchar(255) CHARACTER SET `ascii` NOT NULL;
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer ADD CONSTRAINT `hydra_oauth2_trusted_jwt_bearer_issuer_ibfk_1` FOREIGN KEY (`key_set`, `key_id`) REFERENCES `hydra_jwk` (`sid`, `kid`) ON DELETE CASCADE;
--split

-- hydra_client
ALTER TABLE `hydra_client` ADD COLUMN `nid` char(36);
ALTER TABLE `hydra_client` ADD CONSTRAINT `hydra_client_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_client SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_client MODIFY `nid` char(36) NOT NULL;
--split
CREATE UNIQUE INDEX hydra_client_id_key ON hydra_client (id ASC, nid ASC);
--split



-- hydra_oauth2_access
ALTER TABLE hydra_oauth2_access ADD COLUMN `nid` char(36);
ALTER TABLE hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_nid_fk_idx FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_access SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_oauth2_access MODIFY `nid` char(36) NOT NULL;
--split
ALTER TABLE hydra_oauth2_access DROP FOREIGN KEY `hydra_oauth2_access_client_id_fk`;
ALTER TABLE hydra_oauth2_access ADD CONSTRAINT `hydra_oauth2_access_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE;
--split
DROP INDEX hydra_oauth2_access_requested_at_idx ON hydra_oauth2_access;
DROP INDEX hydra_oauth2_access_request_id_idx ON hydra_oauth2_access;
--split
CREATE INDEX hydra_oauth2_access_requested_at_idx ON hydra_oauth2_access (requested_at, nid);
CREATE INDEX hydra_oauth2_access_client_id_subject_nid_idx ON hydra_oauth2_access (client_id, subject, nid);
CREATE INDEX hydra_oauth2_access_request_id_idx ON hydra_oauth2_access (request_id, nid);
--split



-- hydra_oauth2_authentication_session
ALTER TABLE hydra_oauth2_authentication_session ADD COLUMN `nid` char(36);
ALTER TABLE hydra_oauth2_authentication_session ADD CONSTRAINT hydra_oauth2_authentication_session_nid_fk_idx FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_authentication_session SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_oauth2_authentication_session MODIFY `nid` char(36) NOT NULL;
--split
CREATE INDEX hydra_oauth2_authentication_session_subject_nid_idx ON hydra_oauth2_authentication_session (subject ASC, nid ASC);
--split



-- hydra_oauth2_code
ALTER TABLE `hydra_oauth2_code` ADD COLUMN `nid` char(36);
ALTER TABLE `hydra_oauth2_code` ADD CONSTRAINT `hydra_oauth2_code_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_code SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_oauth2_code MODIFY `nid` char(36) NOT NULL;
--split
ALTER TABLE hydra_oauth2_code DROP FOREIGN KEY `hydra_oauth2_code_client_id_fk`;
ALTER TABLE hydra_oauth2_code ADD CONSTRAINT `hydra_oauth2_code_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE;
--split
DROP INDEX hydra_oauth2_code_request_id_idx ON hydra_oauth2_code;
--split
CREATE INDEX hydra_oauth2_code_request_id_idx ON hydra_oauth2_code (request_id, nid);



-- hydra_oauth2_flow
ALTER TABLE `hydra_oauth2_flow` ADD COLUMN `nid` char(36);
ALTER TABLE `hydra_oauth2_flow` ADD CONSTRAINT `hydra_oauth2_flow_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;
--split
ALTER TABLE hydra_oauth2_flow DROP FOREIGN KEY `hydra_oauth2_flow_client_id_fk`;
ALTER TABLE hydra_oauth2_flow ADD CONSTRAINT `hydra_oauth2_flow_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_flow SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_oauth2_flow MODIFY `nid` char(36) NOT NULL;
--split
DROP INDEX hydra_oauth2_flow_client_id_subject_idx ON hydra_oauth2_flow;
-- DROP INDEX hydra_oauth2_flow_login_session_id_idx ON hydra_oauth2_flow;
DROP INDEX hydra_oauth2_flow_sub_idx ON hydra_oauth2_flow;
DROP INDEX hydra_oauth2_flow_login_verifier_idx ON hydra_oauth2_flow;
DROP INDEX hydra_oauth2_flow_consent_verifier_idx ON hydra_oauth2_flow;
--split
CREATE INDEX hydra_oauth2_flow_client_id_subject_idx ON hydra_oauth2_flow (client_id ASC, nid ASC, subject ASC);
-- CREATE INDEX hydra_oauth2_flow_login_session_id_idx ON hydra_oauth2_flow (login_session_id ASC, nid ASC);
CREATE INDEX hydra_oauth2_flow_sub_idx ON hydra_oauth2_flow (subject ASC, nid ASC);
CREATE UNIQUE INDEX hydra_oauth2_flow_login_verifier_idx ON hydra_oauth2_flow (login_verifier ASC);
CREATE UNIQUE INDEX hydra_oauth2_flow_consent_verifier_idx ON hydra_oauth2_flow (consent_verifier ASC);



-- hydra_oauth2_jti_blacklist
--split
ALTER TABLE `hydra_oauth2_jti_blacklist` ADD COLUMN `nid` char(36);
ALTER TABLE `hydra_oauth2_jti_blacklist` ADD CONSTRAINT `hydra_oauth2_jti_blacklist_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_jti_blacklist SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_oauth2_jti_blacklist MODIFY `nid` char(36) NOT NULL;
--split
DROP INDEX hydra_oauth2_jti_blacklist_expiry ON hydra_oauth2_jti_blacklist;
--split
CREATE INDEX hydra_oauth2_jti_blacklist_expiry ON hydra_oauth2_jti_blacklist (expires_at ASC, nid ASC);
--split
ALTER TABLE hydra_oauth2_jti_blacklist DROP PRIMARY KEY, ADD PRIMARY KEY (signature, nid);
--split



-- hydra_oauth2_logout_request
ALTER TABLE hydra_oauth2_logout_request ADD COLUMN `nid` char(36);
ALTER TABLE hydra_oauth2_logout_request ADD CONSTRAINT hydra_oauth2_logout_request_nid_fk_idx FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;

--split
ALTER TABLE hydra_oauth2_logout_request DROP FOREIGN KEY `hydra_oauth2_logout_request_client_id_fk`;
ALTER TABLE hydra_oauth2_logout_request ADD CONSTRAINT `hydra_oauth2_logout_request_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE;

--split
UPDATE hydra_oauth2_logout_request SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_oauth2_logout_request MODIFY `nid` char(36) NOT NULL;



--split
-- hydra_oauth2_obfuscated_authentication_session
ALTER TABLE hydra_oauth2_obfuscated_authentication_session ADD COLUMN `nid` char(36);
ALTER TABLE hydra_oauth2_obfuscated_authentication_session ADD CONSTRAINT hydra_oauth2_obfuscated_authentication_session_nid_fk_idx FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;
--split
ALTER TABLE hydra_oauth2_obfuscated_authentication_session DROP FOREIGN KEY `hydra_oauth2_obfuscated_authentication_session_client_id_fk`;
ALTER TABLE hydra_oauth2_obfuscated_authentication_session ADD CONSTRAINT `hydra_oauth2_obfuscated_authentication_session_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_obfuscated_authentication_session SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_oauth2_obfuscated_authentication_session MODIFY `nid` char(36) NOT NULL;
--split
ALTER TABLE hydra_oauth2_obfuscated_authentication_session DROP PRIMARY KEY, ADD PRIMARY KEY (subject, client_id, nid);
--split
CREATE UNIQUE INDEX hydra_oauth2_obfuscated_authentication_session_so_nid_idx ON hydra_oauth2_obfuscated_authentication_session (client_id ASC, subject_obfuscated ASC, nid ASC);
--split



-- hydra_oauth2_oidc
ALTER TABLE `hydra_oauth2_oidc` ADD COLUMN `nid` char(36);
ALTER TABLE `hydra_oauth2_oidc` ADD CONSTRAINT `hydra_oauth2_oidc_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_oidc SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_oauth2_oidc MODIFY `nid` char(36) NOT NULL;
--split
ALTER TABLE hydra_oauth2_oidc DROP FOREIGN KEY `hydra_oauth2_oidc_client_id_fk`;
ALTER TABLE hydra_oauth2_oidc ADD CONSTRAINT `hydra_oauth2_oidc_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE;
--split
DROP INDEX hydra_oauth2_oidc_request_id_idx ON hydra_oauth2_oidc;
--split
CREATE INDEX hydra_oauth2_oidc_request_id_idx ON hydra_oauth2_oidc (request_id ASC, nid ASC);



--split
-- hydra_oauth2_pkce
ALTER TABLE hydra_oauth2_pkce ADD COLUMN `nid` char(36);
ALTER TABLE hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_nid_fk_idx FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_pkce SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_oauth2_pkce MODIFY `nid` char(36) NOT NULL;
--split
ALTER TABLE hydra_oauth2_pkce DROP FOREIGN KEY `hydra_oauth2_pkce_client_id_fk`;
ALTER TABLE hydra_oauth2_pkce ADD CONSTRAINT `hydra_oauth2_pkce_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE;
--split
-- DROP INDEX hydra_oauth2_pkce_challenge_id_idx ON hydra_oauth2_pkce;
DROP INDEX hydra_oauth2_pkce_request_id_idx ON hydra_oauth2_pkce;
--split
-- CREATE INDEX hydra_oauth2_pkce_challenge_id_idx ON hydra_oauth2_pkce (challenge_id ASC);
CREATE INDEX hydra_oauth2_pkce_request_id_idx ON hydra_oauth2_pkce (request_id ASC, nid ASC);



--split
-- hydra_oauth2_refresh
ALTER TABLE hydra_oauth2_refresh ADD COLUMN `nid` char(36);
ALTER TABLE hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_nid_fk_idx FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_refresh SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_oauth2_refresh MODIFY `nid` char(36) NOT NULL;
--split
ALTER TABLE hydra_oauth2_refresh DROP FOREIGN KEY `hydra_oauth2_refresh_client_id_fk`;
ALTER TABLE hydra_oauth2_refresh ADD CONSTRAINT `hydra_oauth2_refresh_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE;
--split
-- DROP INDEX hydra_oauth2_refresh_challenge_id_idx ON hydra_oauth2_refresh;
DROP INDEX hydra_oauth2_refresh_client_id_subject_idx ON hydra_oauth2_refresh;
DROP INDEX hydra_oauth2_refresh_request_id_idx ON hydra_oauth2_refresh;
--split
-- CREATE INDEX hydra_oauth2_refresh_challenge_id_idx ON hydra_oauth2_refresh (challenge_id ASC);
CREATE INDEX hydra_oauth2_refresh_client_id_subject_idx ON hydra_oauth2_refresh (client_id ASC, subject ASC);
CREATE INDEX hydra_oauth2_refresh_request_id_idx ON hydra_oauth2_refresh (request_id ASC);



-- hydra_jwk
ALTER TABLE hydra_jwk ADD COLUMN `nid` char(36);
ALTER TABLE hydra_jwk ADD CONSTRAINT hydra_jwk_nid_fk_idx FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_jwk SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_jwk MODIFY `nid` char(36) NOT NULL;
--split
CREATE UNIQUE INDEX hydra_jwk_sid_kid_nid_key ON hydra_jwk (sid ASC, kid ASC, nid ASC);



--split
-- hydra_oauth2_trusted_jwt_bearer_issuer
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer ADD COLUMN `nid` char(36);
--split
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer ADD CONSTRAINT hydra_oauth2_trusted_jwt_bearer_issuer_nid_fk_idx FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE;
--split
UPDATE hydra_oauth2_trusted_jwt_bearer_issuer SET nid = (SELECT id FROM networks LIMIT 1);
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer MODIFY `nid` char(36) NOT NULL;
--split
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer DROP FOREIGN KEY `hydra_oauth2_trusted_jwt_bearer_issuer_ibfk_1`;
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer ADD CONSTRAINT `hydra_oauth2_trusted_jwt_bearer_issuer_ibfk_1` FOREIGN KEY (`key_set`, `key_id`, `nid`) REFERENCES `hydra_jwk` (`sid`, `kid`, `nid`) ON DELETE CASCADE;
--split
DROP INDEX issuer ON hydra_oauth2_trusted_jwt_bearer_issuer;
DROP INDEX hydra_oauth2_trusted_jwt_bearer_issuer_expires_at_idx ON hydra_oauth2_trusted_jwt_bearer_issuer;
--split
CREATE UNIQUE INDEX issuer ON hydra_oauth2_trusted_jwt_bearer_issuer (issuer, subject, key_id, nid);
CREATE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_expires_at_idx ON hydra_oauth2_trusted_jwt_bearer_issuer (expires_at ASC);
CREATE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_nid_idx ON hydra_oauth2_trusted_jwt_bearer_issuer (id, nid);



--split
DROP INDEX hydra_jwk_idx_id_uq ON hydra_jwk;
DROP INDEX hydra_client_idx_id_uq ON hydra_client;
DROP INDEX hydra_oauth2_access_client_id_subject_idx ON hydra_oauth2_access;
-- DROP INDEX hydra_oauth2_authentication_session_subject_idx ON hydra_oauth2_authentication_session;
DROP INDEX hydra_oauth2_flow_cid_idx ON hydra_oauth2_flow;
DROP INDEX hydra_oauth2_obfuscated_authentication_session_so_idx ON hydra_oauth2_obfuscated_authentication_session;

DROP INDEX hydra_oauth2_logout_request_client_id_idx ON hydra_oauth2_logout_request;
DROP INDEX hydra_oauth2_code_client_id_idx ON hydra_oauth2_code;
DROP INDEX hydra_oauth2_access_client_id_idx ON hydra_oauth2_access;
