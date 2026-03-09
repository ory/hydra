-- migrations hash: 36ad8390f65c43551e28df0dcf92b0fdeb823b774eecda791f5979deafce2f6cc6ff57fccdfa41cbaa4403fd4ced8a9dfd7b725d1bb51fd6b0fabdccd51338aa


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

DROP TABLE IF EXISTS `hydra_client`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_client` (
  `id` varchar(255) NOT NULL,
  `client_name` text NOT NULL,
  `client_secret` text NOT NULL,
  `scope` text NOT NULL,
  `owner` text NOT NULL,
  `policy_uri` text NOT NULL,
  `tos_uri` text NOT NULL,
  `client_uri` text NOT NULL,
  `logo_uri` text NOT NULL,
  `client_secret_expires_at` int NOT NULL DEFAULT '0',
  `sector_identifier_uri` text NOT NULL,
  `jwks` text NOT NULL,
  `jwks_uri` text NOT NULL,
  `token_endpoint_auth_method` varchar(25) NOT NULL DEFAULT '',
  `request_object_signing_alg` varchar(10) NOT NULL DEFAULT '',
  `userinfo_signed_response_alg` varchar(10) NOT NULL DEFAULT '',
  `subject_type` varchar(15) NOT NULL DEFAULT '',
  `pk_deprecated` int unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `frontchannel_logout_uri` text NOT NULL,
  `frontchannel_logout_session_required` tinyint(1) NOT NULL DEFAULT '0',
  `backchannel_logout_uri` text NOT NULL,
  `backchannel_logout_session_required` tinyint(1) NOT NULL DEFAULT '0',
  `metadata` text NOT NULL,
  `token_endpoint_auth_signing_alg` varchar(10) NOT NULL DEFAULT '',
  `authorization_code_grant_access_token_lifespan` bigint DEFAULT NULL,
  `authorization_code_grant_id_token_lifespan` bigint DEFAULT NULL,
  `authorization_code_grant_refresh_token_lifespan` bigint DEFAULT NULL,
  `client_credentials_grant_access_token_lifespan` bigint DEFAULT NULL,
  `implicit_grant_access_token_lifespan` bigint DEFAULT NULL,
  `implicit_grant_id_token_lifespan` bigint DEFAULT NULL,
  `jwt_bearer_grant_access_token_lifespan` bigint DEFAULT NULL,
  `password_grant_access_token_lifespan` bigint DEFAULT NULL,
  `password_grant_refresh_token_lifespan` bigint DEFAULT NULL,
  `refresh_token_grant_id_token_lifespan` bigint DEFAULT NULL,
  `refresh_token_grant_access_token_lifespan` bigint DEFAULT NULL,
  `refresh_token_grant_refresh_token_lifespan` bigint DEFAULT NULL,
  `pk` char(36) DEFAULT NULL,
  `registration_access_token_signature` varchar(128) NOT NULL DEFAULT '',
  `nid` char(36) NOT NULL,
  `redirect_uris` json NOT NULL,
  `grant_types` json NOT NULL,
  `response_types` json NOT NULL,
  `audience` json NOT NULL,
  `allowed_cors_origins` json NOT NULL,
  `contacts` json NOT NULL,
  `request_uris` json NOT NULL,
  `post_logout_redirect_uris` json NOT NULL DEFAULT (_utf8mb4'[]'),
  `access_token_strategy` varchar(10) NOT NULL DEFAULT '',
  `skip_consent` tinyint(1) NOT NULL DEFAULT '0',
  `skip_logout_consent` tinyint(1) DEFAULT NULL,
  `device_authorization_grant_id_token_lifespan` bigint DEFAULT NULL,
  `device_authorization_grant_access_token_lifespan` bigint DEFAULT NULL,
  `device_authorization_grant_refresh_token_lifespan` bigint DEFAULT NULL,
  `rotated_secrets` text NOT NULL,
  PRIMARY KEY (`id`,`nid`),
  UNIQUE KEY `hydra_client_id_key` (`id`,`nid`),
  KEY `pk_deprecated` (`pk_deprecated`),
  KEY `hydra_client_nid_fk_idx` (`nid`),
  CONSTRAINT `hydra_client_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_jwk`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_jwk` (
  `sid` varchar(255) NOT NULL,
  `kid` varchar(255) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
  `version` int NOT NULL DEFAULT '0',
  `keydata` text NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `pk_deprecated` int unsigned NOT NULL AUTO_INCREMENT,
  `pk` char(36) NOT NULL,
  `nid` char(36) NOT NULL,
  PRIMARY KEY (`pk`),
  UNIQUE KEY `hydra_jwk_sid_kid_nid_key` (`sid`,`kid`,`nid`),
  KEY `pk_deprecated` (`pk_deprecated`),
  KEY `hydra_jwk_nid_sid_created_at_idx` (`nid`,`sid`,`created_at`),
  KEY `hydra_jwk_nid_sid_kid_created_at_idx` (`nid`,`sid`,`kid`,`created_at`),
  CONSTRAINT `hydra_jwk_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_oauth2_access`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_oauth2_access` (
  `signature` varchar(255) NOT NULL,
  `request_id` varchar(40) NOT NULL DEFAULT '',
  `requested_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `client_id` varchar(255) NOT NULL DEFAULT '',
  `scope` text NOT NULL,
  `granted_scope` text NOT NULL,
  `form_data` text NOT NULL,
  `session_data` text NOT NULL,
  `subject` varchar(255) NOT NULL DEFAULT '',
  `active` tinyint(1) NOT NULL DEFAULT '1',
  `requested_audience` text NOT NULL,
  `granted_audience` text NOT NULL,
  `challenge_id` varchar(40) DEFAULT NULL,
  `nid` char(36) NOT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`signature`),
  KEY `hydra_oauth2_access_challenge_id_idx` (`challenge_id`),
  KEY `hydra_oauth2_access_nid_fk_idx` (`nid`),
  KEY `hydra_oauth2_access_client_id_fk` (`client_id`,`nid`),
  KEY `hydra_oauth2_access_requested_at_idx` (`requested_at`,`nid`),
  KEY `hydra_oauth2_access_client_id_subject_nid_idx` (`client_id`,`subject`,`nid`),
  KEY `hydra_oauth2_access_request_id_idx` (`request_id`,`nid`),
  CONSTRAINT `hydra_oauth2_access_challenge_id_fk` FOREIGN KEY (`challenge_id`) REFERENCES `hydra_oauth2_flow` (`consent_challenge_id`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_access_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_access_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_oauth2_authentication_session`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_oauth2_authentication_session` (
  `id` varchar(40) NOT NULL,
  `authenticated_at` timestamp NULL DEFAULT NULL,
  `subject` varchar(255) NOT NULL,
  `remember` tinyint(1) NOT NULL DEFAULT '0',
  `nid` char(36) NOT NULL,
  `identity_provider_session_id` varchar(40) DEFAULT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `hydra_oauth2_authentication_session_sub_idx` (`subject`),
  KEY `hydra_oauth2_authentication_session_nid_fk_idx` (`nid`),
  KEY `hydra_oauth2_authentication_session_subject_nid_idx` (`subject`,`nid`),
  CONSTRAINT `hydra_oauth2_authentication_session_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_oauth2_code`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_oauth2_code` (
  `signature` varchar(255) NOT NULL,
  `request_id` varchar(40) NOT NULL DEFAULT '',
  `requested_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `client_id` varchar(255) NOT NULL DEFAULT '',
  `scope` text NOT NULL,
  `granted_scope` text NOT NULL,
  `form_data` text NOT NULL,
  `session_data` text NOT NULL,
  `subject` varchar(255) NOT NULL DEFAULT '',
  `active` tinyint(1) NOT NULL DEFAULT '1',
  `requested_audience` text NOT NULL,
  `granted_audience` text NOT NULL,
  `challenge_id` varchar(40) DEFAULT NULL,
  `nid` char(36) NOT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`signature`),
  KEY `hydra_oauth2_code_challenge_id_idx` (`challenge_id`),
  KEY `hydra_oauth2_code_nid_fk_idx` (`nid`),
  KEY `hydra_oauth2_code_client_id_fk` (`client_id`,`nid`),
  CONSTRAINT `hydra_oauth2_code_challenge_id_fk` FOREIGN KEY (`challenge_id`) REFERENCES `hydra_oauth2_flow` (`consent_challenge_id`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_code_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_code_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_oauth2_device_auth_codes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_oauth2_device_auth_codes` (
  `device_code_signature` varchar(255) NOT NULL,
  `user_code_signature` varchar(255) NOT NULL,
  `request_id` varchar(40) NOT NULL,
  `requested_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `client_id` varchar(255) NOT NULL,
  `scope` varchar(1024) NOT NULL,
  `granted_scope` varchar(1024) NOT NULL,
  `form_data` varchar(4096) NOT NULL,
  `session_data` text NOT NULL,
  `subject` varchar(255) NOT NULL DEFAULT '',
  `device_code_active` tinyint(1) NOT NULL DEFAULT '1',
  `user_code_state` smallint NOT NULL DEFAULT '0',
  `requested_audience` varchar(1024) NOT NULL,
  `granted_audience` varchar(1024) NOT NULL,
  `challenge_id` varchar(40) DEFAULT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  `nid` char(36) NOT NULL,
  PRIMARY KEY (`device_code_signature`,`nid`),
  UNIQUE KEY `hydra_oauth2_device_auth_codes_user_code_signature_idx` (`nid`,`user_code_signature`),
  KEY `hydra_oauth2_device_auth_codes_request_id_idx` (`request_id`,`nid`),
  KEY `hydra_oauth2_device_auth_codes_client_id_idx` (`client_id`,`nid`),
  KEY `hydra_oauth2_device_auth_codes_challenge_id_idx` (`challenge_id`),
  CONSTRAINT `hydra_oauth2_device_auth_codes_ibfk_1` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_device_auth_codes_ibfk_2` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `hydra_oauth2_device_auth_codes_ibfk_3` FOREIGN KEY (`challenge_id`) REFERENCES `hydra_oauth2_flow` (`consent_challenge_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_oauth2_flow`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_oauth2_flow` (
  `login_challenge` varchar(40) NOT NULL,
  `login_verifier` varchar(40) DEFAULT NULL,
  `login_csrf` varchar(40) DEFAULT NULL,
  `subject` varchar(255) DEFAULT NULL,
  `request_url` text,
  `login_skip` tinyint(1) DEFAULT NULL,
  `client_id` varchar(255) DEFAULT NULL,
  `requested_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `login_initialized_at` timestamp NULL DEFAULT NULL,
  `oidc_context` json DEFAULT NULL,
  `login_session_id` varchar(40) DEFAULT NULL,
  `state` smallint DEFAULT NULL,
  `login_remember` tinyint(1) DEFAULT NULL,
  `login_remember_for` int DEFAULT NULL,
  `login_error` text,
  `acr` text,
  `login_authenticated_at` timestamp NULL DEFAULT NULL,
  `login_was_used` tinyint(1) DEFAULT NULL,
  `forced_subject_identifier` varchar(255) DEFAULT NULL,
  `context` json DEFAULT NULL,
  `consent_challenge_id` varchar(40) DEFAULT NULL,
  `consent_skip` tinyint(1) DEFAULT NULL,
  `consent_verifier` varchar(40) DEFAULT NULL,
  `consent_csrf` varchar(40) DEFAULT NULL,
  `consent_remember` tinyint(1) DEFAULT NULL,
  `consent_remember_for` int DEFAULT NULL,
  `consent_handled_at` timestamp NULL DEFAULT NULL,
  `consent_error` text,
  `session_access_token` json DEFAULT NULL,
  `session_id_token` json DEFAULT NULL,
  `consent_was_used` tinyint(1) DEFAULT NULL,
  `nid` char(36) NOT NULL,
  `requested_scope` json DEFAULT NULL,
  `requested_at_audience` json DEFAULT NULL,
  `amr` json DEFAULT NULL,
  `granted_scope` json DEFAULT NULL,
  `granted_at_audience` json DEFAULT NULL,
  `login_extend_session_lifespan` tinyint(1) DEFAULT NULL,
  `identity_provider_session_id` varchar(40) DEFAULT NULL,
  `device_challenge_id` varchar(255) DEFAULT NULL,
  `device_code_request_id` varchar(255) DEFAULT NULL,
  `device_verifier` varchar(40) DEFAULT NULL,
  `device_csrf` varchar(40) DEFAULT NULL,
  `device_was_used` tinyint(1) DEFAULT NULL,
  `device_handled_at` timestamp NULL DEFAULT NULL,
  `device_error` varchar(2048) DEFAULT NULL,
  `expires_at` timestamp GENERATED ALWAYS AS (if((`consent_remember_for` > 0),(`requested_at` + interval `consent_remember_for` second),NULL)) VIRTUAL NULL,
  PRIMARY KEY (`login_challenge`),
  UNIQUE KEY `hydra_oauth2_flow_consent_challenge_idx` (`consent_challenge_id`),
  UNIQUE KEY `hydra_oauth2_flow_device_challenge_idx` (`device_challenge_id`),
  KEY `hydra_oauth2_flow_login_session_id_idx` (`login_session_id`),
  KEY `hydra_oauth2_flow_nid_fk_idx` (`nid`),
  KEY `hydra_oauth2_flow_client_id_subject_idx` (`client_id`,`nid`,`subject`),
  KEY `hydra_oauth2_flow_sub_idx` (`subject`,`nid`),
  KEY `hydra_oauth2_flow_previous_consents_idx` (`subject`,`client_id`,`nid`,`consent_skip`,`consent_error`(2),`consent_remember`),
  CONSTRAINT `hydra_oauth2_flow_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_flow_login_session_id_fk` FOREIGN KEY (`login_session_id`) REFERENCES `hydra_oauth2_authentication_session` (`id`) ON DELETE SET NULL,
  CONSTRAINT `hydra_oauth2_flow_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_oauth2_jti_blacklist`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_oauth2_jti_blacklist` (
  `signature` varchar(64) NOT NULL,
  `expires_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `nid` char(36) NOT NULL,
  PRIMARY KEY (`signature`,`nid`),
  KEY `hydra_oauth2_jti_blacklist_nid_fk_idx` (`nid`),
  KEY `hydra_oauth2_jti_blacklist_expiry` (`expires_at`,`nid`),
  CONSTRAINT `hydra_oauth2_jti_blacklist_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_oauth2_logout_request`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_oauth2_logout_request` (
  `challenge` varchar(36) NOT NULL,
  `verifier` varchar(36) NOT NULL,
  `subject` varchar(255) NOT NULL,
  `sid` varchar(36) NOT NULL,
  `client_id` varchar(255) DEFAULT NULL,
  `request_url` text NOT NULL,
  `redir_url` text NOT NULL,
  `was_used` tinyint(1) NOT NULL DEFAULT '0',
  `accepted` tinyint(1) NOT NULL DEFAULT '0',
  `rejected` tinyint(1) NOT NULL DEFAULT '0',
  `rp_initiated` tinyint(1) NOT NULL DEFAULT '0',
  `nid` char(36) NOT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  `requested_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`challenge`),
  UNIQUE KEY `hydra_oauth2_logout_request_veri_idx` (`verifier`),
  KEY `hydra_oauth2_logout_request_nid_fk_idx` (`nid`),
  KEY `hydra_oauth2_logout_request_client_id_fk` (`client_id`,`nid`),
  CONSTRAINT `hydra_oauth2_logout_request_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_logout_request_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_oauth2_obfuscated_authentication_session`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_oauth2_obfuscated_authentication_session` (
  `subject` varchar(255) NOT NULL,
  `client_id` varchar(255) NOT NULL,
  `subject_obfuscated` varchar(255) NOT NULL,
  `nid` char(36) NOT NULL,
  PRIMARY KEY (`subject`,`client_id`,`nid`),
  UNIQUE KEY `hydra_oauth2_obfuscated_authentication_session_so_nid_idx` (`client_id`,`subject_obfuscated`,`nid`),
  KEY `hydra_oauth2_obfuscated_authentication_session_nid_fk_idx` (`nid`),
  KEY `hydra_oauth2_obfuscated_authentication_session_client_id_fk` (`client_id`,`nid`),
  CONSTRAINT `hydra_oauth2_obfuscated_authentication_session_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_obfuscated_authentication_session_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_oauth2_oidc`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_oauth2_oidc` (
  `signature` varchar(255) NOT NULL,
  `request_id` varchar(40) NOT NULL DEFAULT '',
  `requested_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `client_id` varchar(255) NOT NULL DEFAULT '',
  `scope` text NOT NULL,
  `granted_scope` text NOT NULL,
  `form_data` text NOT NULL,
  `session_data` text NOT NULL,
  `subject` varchar(255) NOT NULL DEFAULT '',
  `active` tinyint(1) NOT NULL DEFAULT '1',
  `requested_audience` text NOT NULL,
  `granted_audience` text NOT NULL,
  `challenge_id` varchar(40) DEFAULT NULL,
  `nid` char(36) NOT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`signature`),
  KEY `hydra_oauth2_oidc_client_id_idx` (`client_id`),
  KEY `hydra_oauth2_oidc_challenge_id_idx` (`challenge_id`),
  KEY `hydra_oauth2_oidc_nid_fk_idx` (`nid`),
  KEY `hydra_oauth2_oidc_client_id_fk` (`client_id`,`nid`),
  CONSTRAINT `hydra_oauth2_oidc_challenge_id_fk` FOREIGN KEY (`challenge_id`) REFERENCES `hydra_oauth2_flow` (`consent_challenge_id`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_oidc_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_oidc_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_oauth2_pkce`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_oauth2_pkce` (
  `signature` varchar(255) NOT NULL,
  `request_id` varchar(40) NOT NULL DEFAULT '',
  `requested_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `client_id` varchar(255) NOT NULL DEFAULT '',
  `scope` text NOT NULL,
  `granted_scope` text NOT NULL,
  `form_data` text NOT NULL,
  `session_data` text NOT NULL,
  `subject` varchar(255) NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  `requested_audience` text NOT NULL,
  `granted_audience` text NOT NULL,
  `challenge_id` varchar(40) DEFAULT NULL,
  `nid` char(36) NOT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`signature`),
  KEY `hydra_oauth2_pkce_client_id_idx` (`client_id`),
  KEY `hydra_oauth2_pkce_challenge_id_idx` (`challenge_id`),
  KEY `hydra_oauth2_pkce_nid_fk_idx` (`nid`),
  KEY `hydra_oauth2_pkce_client_id_fk` (`client_id`,`nid`),
  CONSTRAINT `hydra_oauth2_pkce_challenge_id_fk` FOREIGN KEY (`challenge_id`) REFERENCES `hydra_oauth2_flow` (`consent_challenge_id`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_pkce_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_pkce_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_oauth2_refresh`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_oauth2_refresh` (
  `signature` varchar(255) NOT NULL,
  `request_id` varchar(40) NOT NULL DEFAULT '',
  `requested_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `client_id` varchar(255) NOT NULL DEFAULT '',
  `scope` text NOT NULL,
  `granted_scope` text NOT NULL,
  `form_data` text NOT NULL,
  `session_data` text NOT NULL,
  `subject` varchar(255) NOT NULL DEFAULT '',
  `active` tinyint(1) NOT NULL DEFAULT '1',
  `requested_audience` text NOT NULL,
  `granted_audience` text NOT NULL,
  `challenge_id` varchar(40) DEFAULT NULL,
  `nid` char(36) NOT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  `first_used_at` timestamp NULL DEFAULT NULL,
  `access_token_signature` varchar(255) DEFAULT NULL,
  `used_times` int DEFAULT NULL,
  PRIMARY KEY (`signature`),
  KEY `hydra_oauth2_refresh_client_id_idx` (`client_id`),
  KEY `hydra_oauth2_refresh_challenge_id_idx` (`challenge_id`),
  KEY `hydra_oauth2_refresh_client_id_fk` (`client_id`,`nid`),
  KEY `hydra_oauth2_refresh_request_id_idx` (`request_id`),
  KEY `hydra_oauth2_refresh_requested_at_idx` (`nid`,`requested_at`),
  CONSTRAINT `hydra_oauth2_refresh_challenge_id_fk` FOREIGN KEY (`challenge_id`) REFERENCES `hydra_oauth2_flow` (`consent_challenge_id`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_refresh_client_id_fk` FOREIGN KEY (`client_id`, `nid`) REFERENCES `hydra_client` (`id`, `nid`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_refresh_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `hydra_oauth2_trusted_jwt_bearer_issuer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hydra_oauth2_trusted_jwt_bearer_issuer` (
  `id` varchar(36) NOT NULL,
  `issuer` varchar(255) NOT NULL,
  `subject` varchar(255) NOT NULL,
  `scope` text NOT NULL,
  `key_set` varchar(255) NOT NULL,
  `key_id` varchar(255) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `expires_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `nid` char(36) NOT NULL,
  `allow_any_subject` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `hydra_oauth2_trusted_jwt_bearer_issuer_nid_uq_idx` (`nid`,`key_id`,`issuer`,`subject`),
  KEY `hydra_oauth2_trusted_jwt_bearer_issuer_ibfk_1` (`key_set`,`key_id`,`nid`),
  KEY `hydra_oauth2_trusted_jwt_bearer_issuer_expires_at_idx` (`expires_at`),
  KEY `hydra_oauth2_trusted_jwt_bearer_issuer_nid_idx` (`id`,`nid`),
  CONSTRAINT `hydra_oauth2_trusted_jwt_bearer_issuer_ibfk_1` FOREIGN KEY (`key_set`, `key_id`, `nid`) REFERENCES `hydra_jwk` (`sid`, `kid`, `nid`) ON DELETE CASCADE,
  CONSTRAINT `hydra_oauth2_trusted_jwt_bearer_issuer_nid_fk_idx` FOREIGN KEY (`nid`) REFERENCES `networks` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `networks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `networks` (
  `id` char(36) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `schema_migration`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `schema_migration` (
  `version` varchar(48) NOT NULL,
  `version_self` int NOT NULL DEFAULT '0',
  UNIQUE KEY `schema_migration_version_idx` (`version`),
  KEY `schema_migration_version_self_idx` (`version_self`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

