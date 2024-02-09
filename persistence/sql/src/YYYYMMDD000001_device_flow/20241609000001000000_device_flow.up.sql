CREATE TABLE IF NOT EXISTS hydra_oauth2_device_code
(
    signature             VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id            VARCHAR(40) NOT NULL,
    requested_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id             VARCHAR(255) NOT NULL,
    scope                 TEXT NOT NULL,
    granted_scope         TEXT NOT NULL,
    form_data             TEXT NOT NULL,
    session_data          TEXT NOT NULL,
    subject               VARCHAR(255) NOT NULL DEFAULT '',
    active                BOOL NOT NULL DEFAULT true,
    requested_audience    TEXT NULL DEFAULT '',
    granted_audience      TEXT NULL DEFAULT '',
    challenge_id          VARCHAR(40) NULL,
    expires_at            TIMESTAMP    NULL,
    nid                   UUID NULL
);
CREATE INDEX hydra_oauth2_device_code_request_id_idx ON hydra_oauth2_device_code (request_id, nid);
CREATE INDEX hydra_oauth2_device_code_client_id_idx ON hydra_oauth2_device_code (client_id, nid);
CREATE INDEX hydra_oauth2_device_code_challenge_id_idx ON hydra_oauth2_device_code (challenge_id);

CREATE TABLE IF NOT EXISTS hydra_oauth2_user_code
(
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(40) NOT NULL,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          VARCHAR(255) NOT NULL,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL DEFAULT '',
    active             BOOL         NOT NULL DEFAULT true,
    requested_audience TEXT         NULL DEFAULT '',
    granted_audience   TEXT         NULL DEFAULT '',
    challenge_id       VARCHAR(40)  NULL,
    expires_at         TIMESTAMP    NULL,
    nid                UUID         NULL
);
CREATE INDEX hydra_oauth2_user_code_request_id_idx ON hydra_oauth2_user_code (request_id, nid);
CREATE INDEX hydra_oauth2_user_code_client_id_idx ON hydra_oauth2_user_code (client_id, nid);
CREATE INDEX hydra_oauth2_user_code_challenge_id_idx ON hydra_oauth2_user_code (challenge_id);

ALTER TABLE hydra_oauth2_flow ADD COLUMN device_challenge_id VARCHAR(255) NULL;
ALTER TABLE hydra_oauth2_flow ADD COLUMN device_code_request_id VARCHAR(255) NULL;
ALTER TABLE hydra_oauth2_flow ADD COLUMN device_verifier VARCHAR(40) NULL;
ALTER TABLE hydra_oauth2_flow ADD COLUMN device_csrf VARCHAR(40) NULL;
ALTER TABLE hydra_oauth2_flow ADD COLUMN device_user_code_accepted_at TIMESTAMP NULL;
ALTER TABLE hydra_oauth2_flow ADD COLUMN device_was_used BOOLEAN NULL;
ALTER TABLE hydra_oauth2_flow ADD COLUMN device_handled_at TIMESTAMP NULL;
ALTER TABLE hydra_oauth2_flow ADD COLUMN device_error TEXT NULL;

CREATE INDEX hydra_oauth2_flow_device_challenge_idx ON hydra_oauth2_flow (device_challenge_id);

ALTER TABLE hydra_client ADD COLUMN device_authorization_grant_id_token_lifespan BIGINT NULL DEFAULT NULL;
ALTER TABLE hydra_client ADD COLUMN device_authorization_grant_access_token_lifespan BIGINT NULL DEFAULT NULL;
ALTER TABLE hydra_client ADD COLUMN device_authorization_grant_refresh_token_lifespan BIGINT NULL DEFAULT NULL;
