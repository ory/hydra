CREATE TABLE hydra_oauth2_device_grant_request
(
    challenge             VARCHAR(40)  NOT NULL PRIMARY KEY,
    requested_scope       TEXT         NOT NULL,
    verifier              VARCHAR(40)  NOT NULL UNIQUE,
    client_id             VARCHAR(255) NULL,
    request_url           TEXT         NOT NULL,
    requested_audience    VARCHAR(255) NULL     DEFAULT '',
    csrf                  VARCHAR(40)  NOT NULL,
    device_code_signature VARCHAR(255) NULL,
    accepted              BOOL         NOT NULL DEFAULT true,
    accepted_at           TIMESTAMP    NULL,
    nid                   UUID         NULL
);
CREATE INDEX hydra_oauth2_device_grant_request_client_id_idx ON hydra_oauth2_device_grant_request (client_id, nid);
CREATE INDEX hydra_oauth2_device_grant_request_verifier_idx ON hydra_oauth2_device_grant_request (verifier, nid);
CREATE INDEX hydra_oauth2_device_grant_request_challenge_idx ON hydra_oauth2_device_grant_request (challenge, nid);
ALTER TABLE hydra_oauth2_device_grant_request ADD CONSTRAINT hydra_oauth2_device_grant_request_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES hydra_client(id, nid) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_device_grant_request ADD CONSTRAINT hydra_oauth2_device_grant_request_nid_fk_idx FOREIGN KEY (nid) REFERENCES networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;
