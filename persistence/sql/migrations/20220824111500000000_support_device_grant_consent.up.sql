CREATE TABLE hydra_oauth2_device_grant_request
(
    challenge             VARCHAR(40)   NOT NULL PRIMARY KEY,
    requested_scope       TEXT          NOT NULL,
    verifier              VARCHAR(40)   NOT NULL,
    client_id             VARCHAR(255)  NOT NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    requested_at_audience TEXT          NULL     DEFAULT '',
    device_code           VARCHAR(100)  NOT NULL,
    user_code             VARCHAR(10)   NOT NULL
);
CREATE INDEX hydra_oauth2_device_grant_request_client_id_idx ON hydra_oauth2_device_grant_request (client_id);
CREATE UNIQUE INDEX hydra_oauth2_device_grant_request_verifier_idx ON hydra_oauth2_device_grant_request (verifier);
CREATE UNIQUE INDEX hydra_oauth2_device_grant_request_user_code_idx ON hydra_oauth2_device_grant_request (user_code);
CREATE UNIQUE INDEX hydra_oauth2_device_grant_request_device_code_idx ON hydra_oauth2_device_grant_request (device_code);


CREATE TABLE hydra_oauth2_device_grant_request_handled
(
    challenge            VARCHAR(40) NOT NULL PRIMARY KEY REFERENCES hydra_oauth2_device_grant_request (challenge) ON DELETE CASCADE,
    was_used             INTEGER     NOT NULL,
    handled_at           TIMESTAMP   NULL
);
