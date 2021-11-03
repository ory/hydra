-- login_used_at

CREATE TABLE hydra_oauth2_authentication_request_tmp
(
    challenge                 VARCHAR(40)  NOT NULL PRIMARY KEY,
    requested_scope           TEXT         NOT NULL,
    verifier                  VARCHAR(40)  NOT NULL,
    csrf                      VARCHAR(40)  NOT NULL,
    subject                   VARCHAR(255) NOT NULL,
    request_url               TEXT         NOT NULL,
    skip                      INTEGER      NOT NULL,
    client_id                 VARCHAR(255) NOT NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    requested_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    login_initialized_at      TIMESTAMP    NULL,
    oidc_context              TEXT         NOT NULL,
    login_session_id          VARCHAR(40)  NULL REFERENCES hydra_oauth2_authentication_session (id) ON DELETE CASCADE DEFAULT '',
    requested_at_audience     TEXT         NULL DEFAULT '',

    state                     INTEGER      NOT NULL DEFAULT 0,

    remember                  INTEGER      NULL,
    remember_for              INTEGER      NULL,
    error                     TEXT         NULL,
    acr                       TEXT         NULL,
    login_authenticated_at    TIMESTAMP    NULL,
    was_used                  INTEGER      NULL,
    forced_subject_identifier VARCHAR(255) NULL     DEFAULT '',
    context                   TEXT         NULL DEFAULT '{}',
    amr                       TEXT         NULL DEFAULT ''

    CHECK (
        state = 0 OR
        state = 1 OR
        (state = 2 AND (
            remember IS NOT NULL AND
            remember_for IS NOT NULL AND
            error IS NOT NULL AND
            acr IS NOT NULL AND
            was_used IS NOT NULL AND
            context IS NOT NULL AND
            amr IS NOT NULL
        )) OR
        (state = 3 AND (
            remember IS NOT NULL AND
            remember_for IS NOT NULL AND
            error IS NOT NULL AND
            acr IS NOT NULL AND
            was_used IS NOT NULL AND
            context IS NOT NULL AND
            amr IS NOT NULL
        ))
    )
);


INSERT INTO hydra_oauth2_authentication_request_tmp (
    challenge,
    requested_scope,
    verifier,
    csrf,
    subject,
    request_url,
    skip,
    client_id,
    requested_at,
    login_initialized_at,
    oidc_context,
    login_session_id,
    requested_at_audience,

    remember,
    remember_for,
    error,
    acr,
    login_authenticated_at,
    was_used,
    forced_subject_identifier,
    context,
    amr
) SELECT
    hydra_oauth2_authentication_request.challenge,
    hydra_oauth2_authentication_request.requested_scope,
    hydra_oauth2_authentication_request.verifier,
    hydra_oauth2_authentication_request.csrf,
    hydra_oauth2_authentication_request.subject,
    hydra_oauth2_authentication_request.request_url,
    hydra_oauth2_authentication_request.skip,
    hydra_oauth2_authentication_request.client_id,
    hydra_oauth2_authentication_request.requested_at,
    hydra_oauth2_authentication_request.authenticated_at,
    hydra_oauth2_authentication_request.oidc_context,
    hydra_oauth2_authentication_request.login_session_id,
    hydra_oauth2_authentication_request.requested_at_audience,

    hydra_oauth2_authentication_request_handled.remember,
    hydra_oauth2_authentication_request_handled.remember_for,
    hydra_oauth2_authentication_request_handled.error,
    hydra_oauth2_authentication_request_handled.acr,
    hydra_oauth2_authentication_request_handled.authenticated_at,
    hydra_oauth2_authentication_request_handled.was_used,
    hydra_oauth2_authentication_request_handled.forced_subject_identifier,
    hydra_oauth2_authentication_request_handled.context,
    hydra_oauth2_authentication_request_handled.amr
FROM hydra_oauth2_authentication_request
LEFT JOIN hydra_oauth2_authentication_request_handled
ON hydra_oauth2_authentication_request_handled.challenge = hydra_oauth2_authentication_request_handled.challenge;

DROP TABLE hydra_oauth2_authentication_request;
DROP TABLE hydra_oauth2_authentication_request_handled;
ALTER TABLE hydra_oauth2_authentication_request_tmp RENAME TO hydra_oauth2_flow;

CREATE INDEX hydra_oauth2_flow_client_id_idx ON hydra_oauth2_flow (client_id);
CREATE INDEX hydra_oauth2_flow_login_session_id_idx ON hydra_oauth2_flow (login_session_id);
CREATE INDEX hydra_oauth2_flow_subject_idx ON hydra_oauth2_flow (subject);
CREATE UNIQUE INDEX hydra_oauth2_flow_verifier_idx ON hydra_oauth2_flow (verifier);


-- Recreate hydra_oauth2_consent_request in order to rename a foreign key reference

CREATE TABLE hydra_oauth2_consent_request_tmp
(
    challenge                 VARCHAR(40)  NOT NULL PRIMARY KEY,
    verifier                  VARCHAR(40)  NOT NULL,
    client_id                 VARCHAR(255) NOT NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    subject                   VARCHAR(255) NOT NULL,
    request_url               TEXT         NOT NULL,
    skip                      INTEGER      NOT NULL,
    requested_scope           TEXT         NOT NULL,
    csrf                      VARCHAR(40)  NOT NULL,
    authenticated_at          TIMESTAMP    NULL,
    requested_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    oidc_context              TEXT         NOT NULL,
    forced_subject_identifier VARCHAR(255) NULL     DEFAULT '',
    login_session_id          VARCHAR(40)  NULL REFERENCES hydra_oauth2_authentication_session (id) ON DELETE SET NULL,
    login_challenge           VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (challenge) ON DELETE SET NULL,
    requested_at_audience     TEXT         NULL     DEFAULT '',
    acr                       TEXT         NULL     DEFAULT '',
    amr                       TEXT         NOT NULL DEFAULT '',
    context                   TEXT         NOT NULL DEFAULT '{}'
);

INSERT INTO hydra_oauth2_consent_request_tmp (
     challenge,
     verifier,
     client_id,
     subject,
     request_url,
     skip,
     requested_scope,
     csrf,
     authenticated_at,
     requested_at,
     oidc_context,
     forced_subject_identifier,
     login_session_id,
     login_challenge,
     requested_at_audience,
     acr,
     amr,
     context
) SELECT
     challenge,
     verifier,
     client_id,
     subject,
     request_url,
     skip,
     requested_scope,
     csrf,
     authenticated_at,
     requested_at,
     oidc_context,
     forced_subject_identifier,
     login_session_id,
     login_challenge,
     requested_at_audience,
     acr,
     amr,
     context
FROM hydra_oauth2_consent_request;

DROP TABLE hydra_oauth2_consent_request;
ALTER TABLE hydra_oauth2_consent_request_tmp RENAME TO hydra_oauth2_consent_request;

CREATE INDEX hydra_oauth2_consent_request_client_id_idx ON hydra_oauth2_consent_request (client_id);
CREATE INDEX hydra_oauth2_consent_request_subject_idx ON hydra_oauth2_consent_request (subject);
CREATE INDEX hydra_oauth2_consent_request_login_session_id_idx ON hydra_oauth2_consent_request (login_session_id);
CREATE INDEX hydra_oauth2_consent_request_login_challenge_idx ON hydra_oauth2_consent_request (login_challenge);
