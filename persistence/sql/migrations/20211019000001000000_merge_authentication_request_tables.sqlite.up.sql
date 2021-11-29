CREATE TABLE hydra_oauth2_flow
(
    login_challenge           VARCHAR(40)  NOT NULL PRIMARY KEY,
    requested_scope           TEXT         NOT NULL,
    login_verifier            VARCHAR(40)  NOT NULL,
    login_csrf                VARCHAR(40)  NOT NULL,
    subject                   VARCHAR(255) NOT NULL,
    request_url               TEXT         NOT NULL,
    login_skip                INTEGER      NOT NULL,
    client_id                 VARCHAR(255) NOT NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    requested_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    login_initialized_at      TIMESTAMP    NULL,
    oidc_context              TEXT         NOT NULL,
    login_session_id          VARCHAR(40)  NULL REFERENCES hydra_oauth2_authentication_session (id) ON DELETE CASCADE DEFAULT '',
    requested_at_audience     TEXT         NULL DEFAULT '',

    state                     INTEGER      NOT NULL DEFAULT 0,

    login_remember            INTEGER      NULL,
    login_remember_for        INTEGER      NULL,
    login_error               TEXT         NULL,
    acr                       TEXT         NULL,
    login_authenticated_at    TIMESTAMP    NULL,
    login_was_used            INTEGER      NULL,
    forced_subject_identifier VARCHAR(255) NULL DEFAULT '',
    context                   TEXT         NULL DEFAULT '{}',
    amr                       TEXT         NULL DEFAULT '',

    consent_challenge_id      VARCHAR(40)  NULL,
    consent_skip              INTEGER      NOT NULL,
    consent_verifier          VARCHAR(40)  NOT NULL,
    consent_csrf              VARCHAR(40)  NOT NULL,

    granted_scope             TEXT        NOT NULL,
    granted_at_audience       TEXT        NULL DEFAULT '',
    consent_remember          INTEGER     NOT NULL,
    consent_remember_for      INTEGER     NOT NULL,
    consent_handled_at        TIMESTAMP   NULL,
    consent_was_used          INTEGER     NOT NULL,
    consent_error             TEXT        NOT NULL,
    session_id_token          TEXT        NOT NULL,
    session_access_token      TEXT        NOT NULL,

    CHECK (
        state = 128 OR
        state = 129 OR
        state = 1 OR
        (state = 2 AND (
            login_remember IS NOT NULL AND
            login_remember_for IS NOT NULL AND
            login_error IS NOT NULL AND
            acr IS NOT NULL AND
            login_was_used IS NOT NULL AND
            context IS NOT NULL AND
            amr IS NOT NULL
        )) OR
        (state = 3 AND (
            login_remember IS NOT NULL AND
            login_remember_for IS NOT NULL AND
            login_error IS NOT NULL AND
            acr IS NOT NULL AND
            login_was_used IS NOT NULL AND
            context IS NOT NULL AND
            amr IS NOT NULL
        )) OR
        (state = 4 AND (
            login_remember IS NOT NULL AND
            login_remember_for IS NOT NULL AND
            login_error IS NOT NULL AND
            acr IS NOT NULL AND
            login_was_used IS NOT NULL AND
            context IS NOT NULL AND
            amr IS NOT NULL AND

            consent_challenge_id IS NOT NULL AND
            consent_verifier IS NOT NULL AND
            consent_skip IS NOT NULL AND
            consent_csrf IS NOT NULL
        )) OR
        (state = 5 AND (
            login_remember IS NOT NULL AND
            login_remember_for IS NOT NULL AND
            login_error IS NOT NULL AND
            acr IS NOT NULL AND
            login_was_used IS NOT NULL AND
            context IS NOT NULL AND
            amr IS NOT NULL AND

            consent_challenge_id IS NOT NULL AND
            consent_verifier IS NOT NULL AND
            consent_skip IS NOT NULL AND
            consent_csrf IS NOT NULL AND

            granted_scope IS NOT NULL AND
            consent_remember IS NOT NULL AND
            consent_remember_for IS NOT NULL AND
            consent_error IS NOT NULL AND
            session_access_token IS NOT NULL AND
            session_id_token IS NOT NULL AND
            consent_was_used IS NOT NULL
        ))
    )
);

INSERT INTO hydra_oauth2_flow (
    login_challenge,
    requested_scope,
    login_verifier,
    login_csrf,
    subject,
    request_url,
    login_skip,
    client_id,
    requested_at,
    login_initialized_at,
    oidc_context,
    login_session_id,
    requested_at_audience,

    login_remember,
    login_remember_for,
    login_error,
    acr,
    login_authenticated_at,
    login_was_used,
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

CREATE INDEX hydra_oauth2_flow_client_id_idx ON hydra_oauth2_flow (client_id);
CREATE INDEX hydra_oauth2_flow_login_session_id_idx ON hydra_oauth2_flow (login_session_id);
CREATE INDEX hydra_oauth2_flow_subject_idx ON hydra_oauth2_flow (subject);
CREATE UNIQUE INDEX hydra_oauth2_flow_consent_challenge_id_idx ON hydra_oauth2_flow (consent_challenge_id);
CREATE UNIQUE INDEX hydra_oauth2_flow_login_verifier_idx ON hydra_oauth2_flow (login_verifier);
CREATE INDEX hydra_oauth2_flow_consent_verifier_idx ON hydra_oauth2_flow (consent_verifier);


DROP TABLE hydra_oauth2_authentication_request;
DROP TABLE hydra_oauth2_authentication_request_handled;
DROP TABLE hydra_oauth2_consent_request;
DROP TABLE hydra_oauth2_consent_request_handled;


DROP TABLE hydra_oauth2_code;
CREATE TABLE hydra_oauth2_code
(
    signature          VARCHAR(255) NOT NULL,
    request_id         VARCHAR(40)  NOT NULL,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          VARCHAR(255) NOT NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL DEFAULT '',
    active             INTEGER      NOT NULL DEFAULT true,
    requested_audience TEXT         NULL     DEFAULT '',
    granted_audience   TEXT         NULL     DEFAULT '',
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE
);

CREATE INDEX hydra_oauth2_code_client_id_idx ON hydra_oauth2_code (client_id);
CREATE INDEX hydra_oauth2_code_challenge_id_idx ON hydra_oauth2_code (challenge_id);
CREATE INDEX hydra_oauth2_code_request_id_idx ON hydra_oauth2_code (request_id);


DROP TABLE hydra_oauth2_oidc;
CREATE TABLE hydra_oauth2_oidc
(
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(40)  NOT NULL,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          VARCHAR(255) NOT NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL DEFAULT '',
    active             INTEGER      NOT NULL DEFAULT true,
    requested_audience TEXT         NULL     DEFAULT '',
    granted_audience   TEXT         NULL     DEFAULT '',
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE
);

CREATE INDEX hydra_oauth2_oidc_client_id_idx ON hydra_oauth2_oidc (client_id);
CREATE INDEX hydra_oauth2_oidc_challenge_id_idx ON hydra_oauth2_oidc (challenge_id);
CREATE INDEX hydra_oauth2_oidc_request_id_idx ON hydra_oauth2_oidc (request_id);


DROP TABLE hydra_oauth2_pkce;
CREATE TABLE hydra_oauth2_pkce
(
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(40)  NOT NULL,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          VARCHAR(255) NOT NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL,
    active             INTEGER      NOT NULL DEFAULT true,
    requested_audience TEXT         NULL     DEFAULT '',
    granted_audience   TEXT         NULL     DEFAULT '',
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE
);

CREATE INDEX hydra_oauth2_pkce_client_id_idx ON hydra_oauth2_pkce (client_id);
CREATE INDEX hydra_oauth2_pkce_challenge_id_idx ON hydra_oauth2_pkce (challenge_id);
CREATE INDEX hydra_oauth2_pkce_request_id_idx ON hydra_oauth2_pkce (request_id);


DROP TABLE "hydra_oauth2_access";
CREATE TABLE IF NOT EXISTS "hydra_oauth2_access"
(
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(40)  NOT NULL,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          VARCHAR(255) NOT NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL DEFAULT '',
    active             INTEGER      NOT NULL DEFAULT true,
    requested_audience TEXT         NULL     DEFAULT '',
    granted_audience   TEXT         NULL     DEFAULT '',
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE
);

CREATE INDEX hydra_oauth2_access_requested_at_idx ON hydra_oauth2_access (requested_at);
CREATE INDEX hydra_oauth2_access_client_id_idx ON hydra_oauth2_access (client_id);
CREATE INDEX hydra_oauth2_access_challenge_id_idx ON hydra_oauth2_access (challenge_id);
CREATE INDEX hydra_oauth2_access_client_id_subject_idx ON hydra_oauth2_access (client_id, subject);
CREATE INDEX hydra_oauth2_access_request_id_idx ON hydra_oauth2_access (request_id);


DROP TABLE hydra_oauth2_refresh;
CREATE TABLE hydra_oauth2_refresh
(
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(40)  NOT NULL,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          VARCHAR(255) NOT NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL DEFAULT '',
    active             INTEGER      NOT NULL DEFAULT true,
    requested_audience TEXT         NULL     DEFAULT '',
    granted_audience   TEXT         NULL     DEFAULT '',
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE
);

CREATE INDEX hydra_oauth2_refresh_client_id_idx ON hydra_oauth2_refresh (client_id);
CREATE INDEX hydra_oauth2_refresh_challenge_id_idx ON hydra_oauth2_refresh (challenge_id);
CREATE INDEX hydra_oauth2_refresh_client_id_subject_idx ON hydra_oauth2_refresh (client_id, subject);
CREATE INDEX hydra_oauth2_refresh_request_id_idx ON hydra_oauth2_refresh (request_id);
