CREATE TABLE hydra_oauth2_device_grant_request
(
    challenge             VARCHAR(40)  NOT NULL PRIMARY KEY,
    requested_scope       TEXT         NOT NULL,
    verifier              VARCHAR(40)  NOT NULL UNIQUE,
    client_id             VARCHAR(255) NULL REFERENCES hydra_client (pk) ON DELETE CASCADE,
    request_url           TEXT         NOT NULL,
    requested_audience    VARCHAR(255) NULL     DEFAULT '',
    csrf                  VARCHAR(40)  NOT NULL,
    device_request_id     VARCHAR(40)  NULL,
    accepted              BOOL         NOT NULL DEFAULT true,
    accepted_at           TIMESTAMP    NULL
);
CREATE INDEX hydra_oauth2_device_grant_request_client_id_idx ON hydra_oauth2_device_grant_request (client_id);
CREATE INDEX hydra_oauth2_device_grant_request_verifier_idx ON hydra_oauth2_device_grant_request (verifier);
CREATE INDEX hydra_oauth2_device_grant_request_challenge_idx ON hydra_oauth2_device_grant_request (challenge);