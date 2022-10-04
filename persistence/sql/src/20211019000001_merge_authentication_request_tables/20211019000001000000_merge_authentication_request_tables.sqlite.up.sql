CREATE TABLE hydra_oauth2_flow
(
    login_challenge           VARCHAR(40)  NOT NULL PRIMARY KEY,
    requested_scope           TEXT         NOT NULL DEFAULT '[]',
    login_verifier            VARCHAR(40)  NOT NULL,
    login_csrf                VARCHAR(40)  NOT NULL,
    subject                   VARCHAR(255) NOT NULL,
    request_url               TEXT         NOT NULL,
    login_skip                INTEGER      NOT NULL,
    client_id                 VARCHAR(255) NOT NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    requested_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    login_initialized_at      TIMESTAMP    NULL DEFAULT NULL,
    oidc_context              jsonb         NOT NULL DEFAULT '{}',
    login_session_id          VARCHAR(40)  NULL REFERENCES hydra_oauth2_authentication_session (id) ON DELETE CASCADE DEFAULT '',
    requested_at_audience     text         NULL DEFAULT '[]',

    state                     INTEGER      NOT NULL,

    login_remember            INTEGER      NOT NULL DEFAULT false,
    login_remember_for        INTEGER      NOT NULL,
    login_error               TEXT         NULL,
    acr                       TEXT          NOT NULL DEFAULT '',
    login_authenticated_at    TIMESTAMP    NULL DEFAULT NULL,
    login_was_used            INTEGER      NOT NULL DEFAULT false,
    forced_subject_identifier VARCHAR(255) NOT NULL DEFAULT '',
    context                   jsonb         NOT NULL DEFAULT '{}',
    amr                       text         NOT NULL DEFAULT '[]',

    consent_challenge_id      VARCHAR(40)  NULL,
    consent_skip              INTEGER      NOT NULL DEFAULT false,
    consent_verifier          VARCHAR(40)  NULL,
    consent_csrf              VARCHAR(40)  NULL,

    granted_scope             text        NOT NULL DEFAULT '[]',
    granted_at_audience       text        NOT NULL DEFAULT '[]',
    consent_remember          INTEGER     NOT NULL DEFAULT 0,
    consent_remember_for      INTEGER     NULL,
    consent_handled_at        TIMESTAMP   NULL,
    consent_was_used          INTEGER     NOT NULL DEFAULT false,
    consent_error             TEXT        NULL,
    session_id_token          jsonb        NULL DEFAULT '{}',
    session_access_token      jsonb        NULL DEFAULT '{}'

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
            consent_csrf IS NOT NULL
        )) OR
        (state = 6 AND (
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
--split
INSERT INTO hydra_oauth2_flow (
    state,
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
    amr,

    consent_challenge_id,
    consent_verifier,
    consent_skip,
    consent_csrf,

    granted_scope,
    consent_remember,
    consent_remember_for,
    consent_error,
    session_access_token,
    session_id_token,
    consent_was_used,
    granted_at_audience,
    consent_handled_at
) SELECT
    case
        when hydra_oauth2_authentication_request_handled.error IS NOT NULL then 128
        when hydra_oauth2_consent_request_handled.error IS NOT NULL then 129
        when hydra_oauth2_consent_request_handled.was_used = true then 6
        when hydra_oauth2_consent_request_handled.challenge IS NOT NULL then 5
        when hydra_oauth2_consent_request.challenge IS NOT NULL then 4
        when hydra_oauth2_authentication_request_handled.was_used = true then 3
        when hydra_oauth2_authentication_request_handled.challenge IS NOT NULL then 2
        else 1
    end,
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
    coalesce(hydra_oauth2_authentication_request.oidc_context, '{}'),
    hydra_oauth2_authentication_request.login_session_id,
    hydra_oauth2_authentication_request.requested_at_audience,

    coalesce(hydra_oauth2_authentication_request_handled.remember, false),
    coalesce(hydra_oauth2_authentication_request_handled.remember_for, 0),
    hydra_oauth2_authentication_request_handled.error,
    coalesce(hydra_oauth2_authentication_request_handled.acr, ''),
    hydra_oauth2_authentication_request_handled.authenticated_at,
    coalesce(hydra_oauth2_authentication_request_handled.was_used, false),
    coalesce(hydra_oauth2_consent_request.forced_subject_identifier, hydra_oauth2_authentication_request_handled.forced_subject_identifier, ''),
    coalesce(hydra_oauth2_authentication_request_handled.context, '{}'),
    coalesce(hydra_oauth2_authentication_request_handled.amr, ''),

    hydra_oauth2_consent_request.challenge,
    hydra_oauth2_consent_request.verifier,
    coalesce(hydra_oauth2_consent_request.skip, false),
    hydra_oauth2_consent_request.csrf,

    coalesce(hydra_oauth2_consent_request_handled.granted_scope, '[]'),
    coalesce(hydra_oauth2_consent_request_handled.remember, false),
    coalesce(hydra_oauth2_consent_request_handled.remember_for, 0),
    hydra_oauth2_consent_request_handled.error,
    coalesce(hydra_oauth2_consent_request_handled.session_access_token, '{}'),
    coalesce(hydra_oauth2_consent_request_handled.session_id_token, '{}'),
    coalesce(hydra_oauth2_consent_request_handled.was_used, false),
    coalesce(hydra_oauth2_consent_request_handled.granted_at_audience, '[]'),
    hydra_oauth2_consent_request_handled.handled_at
FROM hydra_oauth2_authentication_request
LEFT JOIN hydra_oauth2_authentication_request_handled
ON hydra_oauth2_authentication_request.challenge = hydra_oauth2_authentication_request_handled.challenge
LEFT JOIN hydra_oauth2_consent_request
ON hydra_oauth2_authentication_request.challenge = hydra_oauth2_consent_request.login_challenge
LEFT JOIN hydra_oauth2_consent_request_handled
ON hydra_oauth2_consent_request.challenge = hydra_oauth2_consent_request_handled.challenge;

--split
CREATE INDEX hydra_oauth2_flow_client_id_idx ON hydra_oauth2_flow (client_id);
--split
CREATE INDEX hydra_oauth2_flow_login_session_id_idx ON hydra_oauth2_flow (login_session_id);
--split
CREATE INDEX hydra_oauth2_flow_subject_idx ON hydra_oauth2_flow (subject);
--split
CREATE UNIQUE INDEX hydra_oauth2_flow_consent_challenge_id_idx ON hydra_oauth2_flow (consent_challenge_id);
--split
CREATE UNIQUE INDEX hydra_oauth2_flow_login_verifier_idx ON hydra_oauth2_flow (login_verifier);
--split
CREATE INDEX hydra_oauth2_flow_consent_verifier_idx ON hydra_oauth2_flow (consent_verifier);

--split
DROP TABLE hydra_oauth2_authentication_request;
--split
DROP TABLE hydra_oauth2_authentication_request_handled;
--split
DROP TABLE hydra_oauth2_consent_request;
--split
DROP TABLE hydra_oauth2_consent_request_handled;


--split
DROP TABLE hydra_oauth2_code;
--split
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

--split
CREATE INDEX hydra_oauth2_code_client_id_idx ON hydra_oauth2_code (client_id);
--split
CREATE INDEX hydra_oauth2_code_challenge_id_idx ON hydra_oauth2_code (challenge_id);
--split
CREATE INDEX hydra_oauth2_code_request_id_idx ON hydra_oauth2_code (request_id);


--split
DROP TABLE hydra_oauth2_oidc;
--split
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

--split
CREATE INDEX hydra_oauth2_oidc_client_id_idx ON hydra_oauth2_oidc (client_id);
--split
CREATE INDEX hydra_oauth2_oidc_challenge_id_idx ON hydra_oauth2_oidc (challenge_id);
--split
CREATE INDEX hydra_oauth2_oidc_request_id_idx ON hydra_oauth2_oidc (request_id);


--split
DROP TABLE hydra_oauth2_pkce;
--split
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

--split
CREATE INDEX hydra_oauth2_pkce_client_id_idx ON hydra_oauth2_pkce (client_id);
--split
CREATE INDEX hydra_oauth2_pkce_challenge_id_idx ON hydra_oauth2_pkce (challenge_id);
--split
CREATE INDEX hydra_oauth2_pkce_request_id_idx ON hydra_oauth2_pkce (request_id);


--split
DROP TABLE "hydra_oauth2_access";
--split
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

--split
CREATE INDEX hydra_oauth2_access_requested_at_idx ON hydra_oauth2_access (requested_at);
--split
CREATE INDEX hydra_oauth2_access_client_id_idx ON hydra_oauth2_access (client_id);
--split
CREATE INDEX hydra_oauth2_access_challenge_id_idx ON hydra_oauth2_access (challenge_id);
--split
CREATE INDEX hydra_oauth2_access_client_id_subject_idx ON hydra_oauth2_access (client_id, subject);
--split
CREATE INDEX hydra_oauth2_access_request_id_idx ON hydra_oauth2_access (request_id);


--split
DROP TABLE hydra_oauth2_refresh;
--split
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

--split
CREATE INDEX hydra_oauth2_refresh_client_id_idx ON hydra_oauth2_refresh (client_id);
--split
CREATE INDEX hydra_oauth2_refresh_challenge_id_idx ON hydra_oauth2_refresh (challenge_id);
--split
CREATE INDEX hydra_oauth2_refresh_client_id_subject_idx ON hydra_oauth2_refresh (client_id, subject);
--split
CREATE INDEX hydra_oauth2_refresh_request_id_idx ON hydra_oauth2_refresh (request_id);
