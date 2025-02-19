CREATE TABLE IF NOT EXISTS hydra_oauth2_device_auth_codes
(
  device_code_signature VARCHAR(255)  NOT NULL,
  user_code_signature   VARCHAR(255)  NOT NULL,
  request_id            VARCHAR(40)   NOT NULL,
  requested_at          TIMESTAMP     NOT NULL DEFAULT NOW(),
  client_id             VARCHAR(255)  NOT NULL,
  scope                 VARCHAR(1024) NOT NULL,
  granted_scope         VARCHAR(1024) NOT NULL,
  form_data             VARCHAR(4096) NOT NULL,
  session_data          TEXT NOT NULL,
  subject               VARCHAR(255)  NOT NULL DEFAULT '',
  device_code_active    BOOL          NOT NULL DEFAULT true,
  user_code_state       SMALLINT      NOT NULL DEFAULT 0,
  requested_audience    VARCHAR(1024) NOT NULL,
  granted_audience      VARCHAR(1024) NOT NULL,
  challenge_id          VARCHAR(40)   NULL,
  expires_at            TIMESTAMP     NULL,
  nid                   CHAR(36)      NOT NULL,

  FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE,
  FOREIGN KEY (nid) REFERENCES networks (id) ON UPDATE RESTRICT ON DELETE CASCADE,
  FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE,
  PRIMARY KEY (device_code_signature, nid)
);

CREATE INDEX hydra_oauth2_device_auth_codes_request_id_idx ON hydra_oauth2_device_auth_codes (request_id, nid);
CREATE INDEX hydra_oauth2_device_auth_codes_client_id_idx ON hydra_oauth2_device_auth_codes (client_id, nid);
CREATE INDEX hydra_oauth2_device_auth_codes_challenge_id_idx ON hydra_oauth2_device_auth_codes (challenge_id);
CREATE UNIQUE INDEX hydra_oauth2_device_auth_codes_user_code_signature_idx ON hydra_oauth2_device_auth_codes (nid, user_code_signature);

ALTER TABLE hydra_oauth2_flow
  ADD COLUMN device_challenge_id VARCHAR(255) NULL;
ALTER TABLE hydra_oauth2_flow
  ADD COLUMN device_code_request_id VARCHAR(255) NULL;
ALTER TABLE hydra_oauth2_flow
  ADD COLUMN device_verifier VARCHAR(40) NULL;
ALTER TABLE hydra_oauth2_flow
  ADD COLUMN device_csrf VARCHAR(40) NULL;
ALTER TABLE hydra_oauth2_flow
  ADD COLUMN device_was_used BOOL NULL;
ALTER TABLE hydra_oauth2_flow
  ADD COLUMN device_handled_at TIMESTAMP NULL;
ALTER TABLE hydra_oauth2_flow
  ADD COLUMN device_error VARCHAR(2048) NULL;

CREATE UNIQUE INDEX hydra_oauth2_flow_device_challenge_idx ON hydra_oauth2_flow (device_challenge_id);

ALTER TABLE hydra_client
  ADD COLUMN device_authorization_grant_id_token_lifespan BIGINT NULL DEFAULT NULL;
ALTER TABLE hydra_client
  ADD COLUMN device_authorization_grant_access_token_lifespan BIGINT NULL DEFAULT NULL;
ALTER TABLE hydra_client
  ADD COLUMN device_authorization_grant_refresh_token_lifespan BIGINT NULL DEFAULT NULL;
