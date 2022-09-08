CREATE TABLE IF NOT EXISTS hydra_oauth2_device_code 
(
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(128) NOT NULL UNIQUE,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          UUID         NOT NULL REFERENCES hydra_client (pk) ON DELETE CASCADE,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL ,
    active             BOOL         NOT NULL DEFAULT true,
    requested_audience TEXT         NULL     ,
    granted_audience   TEXT         NULL     ,
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE,
    nid                VARCHAR(36)
);
CREATE INDEX hydra_oauth2_device_code_requested_at_idx ON hydra_oauth2_device_code (requested_at);
CREATE INDEX hydra_oauth2_device_code_client_id_idx ON hydra_oauth2_device_code (client_id);
CREATE INDEX hydra_oauth2_device_code_challenge_id_idx ON hydra_oauth2_device_code (challenge_id);
CREATE INDEX hydra_oauth2_device_code_client_id_subject_idx ON hydra_oauth2_device_code (client_id, subject);


CREATE TABLE IF NOT EXISTS hydra_oauth2_user_code 
(
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(128) NOT NULL UNIQUE,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          UUID         NOT NULL REFERENCES hydra_client (pk) ON DELETE CASCADE,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL DEFAULT '',
    active             BOOL         NOT NULL DEFAULT true,
    requested_audience TEXT         NULL     ,
    granted_audience   TEXT         NULL     ,
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (login_challenge) ON DELETE CASCADE,
    nid                VARCHAR(36)
);
CREATE INDEX hydra_oauth2_user_code_requested_at_idx ON hydra_oauth2_user_code (requested_at);
CREATE INDEX hydra_oauth2_user_code_client_id_idx ON hydra_oauth2_user_code (client_id);
CREATE INDEX hydra_oauth2_user_code_challenge_id_idx ON hydra_oauth2_user_code (challenge_id);
CREATE INDEX hydra_oauth2_user_code_client_id_subject_idx ON hydra_oauth2_user_code (client_id, subject);