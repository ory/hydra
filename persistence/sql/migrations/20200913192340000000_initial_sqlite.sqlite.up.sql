CREATE TABLE IF NOT EXISTS hydra_client
(
    id                                   VARCHAR(255) NOT NULL,
    client_name                          TEXT         NOT NULL,
    client_secret                        TEXT         NOT NULL,
    redirect_uris                        TEXT         NOT NULL,
    grant_types                          TEXT         NOT NULL,
    response_types                       TEXT         NOT NULL,
    scope                                TEXT         NOT NULL,
    owner                                TEXT         NOT NULL,
    policy_uri                           TEXT         NOT NULL,
    tos_uri                              TEXT         NOT NULL,
    client_uri                           TEXT         NOT NULL,
    logo_uri                             TEXT         NOT NULL,
    contacts                             TEXT         NOT NULL,
    client_secret_expires_at             INTEGER      NOT NULL DEFAULT 0,
    sector_identifier_uri                TEXT         NOT NULL,
    jwks                                 TEXT         NOT NULL,
    jwks_uri                             TEXT         NOT NULL,
    request_uris                         TEXT         NOT NULL,
    token_endpoint_auth_method           VARCHAR(25)  NOT NULL DEFAULT '',
    request_object_signing_alg           VARCHAR(10)  NOT NULL DEFAULT '',
    userinfo_signed_response_alg         VARCHAR(10)  NOT NULL DEFAULT '',
    subject_type                         VARCHAR(15)  NOT NULL DEFAULT '',
    allowed_cors_origins                 TEXT         NOT NULL,
    pk                                   INTEGER PRIMARY KEY,
    audience                             TEXT         NOT NULL,
    created_at                           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    frontchannel_logout_uri              TEXT         NOT NULL DEFAULT '',
    frontchannel_logout_session_required INTEGER      NOT NULL DEFAULT false,
    post_logout_redirect_uris            TEXT         NOT NULL DEFAULT '',
    backchannel_logout_uri               TEXT         NOT NULL DEFAULT '',
    backchannel_logout_session_required  INTEGER      NOT NULL DEFAULT false,
    metadata                             TEXT         NOT NULL DEFAULT '{}',
    token_endpoint_auth_signing_alg      VARCHAR(10)  NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX hydra_client_id_idx ON hydra_client (id);

CREATE TABLE IF NOT EXISTS hydra_jwk
(
    sid        VARCHAR(255) NOT NULL,
    kid        VARCHAR(255) NOT NULL,
    version    INTEGER      NOT NULL DEFAULT 0,
    keydata    TEXT         NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    pk         INTEGER PRIMARY KEY
);

CREATE UNIQUE INDEX hydra_jwk_sid_kid_key ON hydra_jwk (sid, kid);

CREATE TABLE hydra_oauth2_authentication_session
(
    id               VARCHAR(40)  NOT NULL PRIMARY KEY,
    authenticated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    subject          VARCHAR(255) NOT NULL,
    remember         INTEGER      NOT NULL DEFAULT false
);

CREATE INDEX hydra_oauth2_authentication_session_subject_idx ON hydra_oauth2_authentication_session (subject);

CREATE TABLE hydra_oauth2_authentication_request
(
    challenge             VARCHAR(40)  NOT NULL PRIMARY KEY,
    requested_scope       TEXT         NOT NULL,
    verifier              VARCHAR(40)  NOT NULL,
    csrf                  VARCHAR(40)  NOT NULL,
    subject               VARCHAR(255) NOT NULL,
    request_url           TEXT         NOT NULL,
    skip                  INTEGER      NOT NULL,
    client_id             VARCHAR(255) NOT NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    requested_at          TIMESTAMP    NOT NULL                                                                               DEFAULT CURRENT_TIMESTAMP,
    authenticated_at      TIMESTAMP    NULL,
    oidc_context          TEXT         NOT NULL,
    login_session_id      VARCHAR(40)  NULL REFERENCES hydra_oauth2_authentication_session (id) ON DELETE CASCADE DEFAULT '',
    requested_at_audience TEXT         NULL                                                                                   DEFAULT ''
);

CREATE INDEX hydra_oauth2_authentication_request_client_id_idx ON hydra_oauth2_authentication_request (client_id);
CREATE INDEX hydra_oauth2_authentication_request_login_session_id_idx ON hydra_oauth2_authentication_request (login_session_id);
CREATE INDEX hydra_oauth2_authentication_request_subject_idx ON hydra_oauth2_authentication_request (subject);
CREATE UNIQUE INDEX hydra_oauth2_authentication_request_verifier_idx ON hydra_oauth2_authentication_request (verifier);

CREATE TABLE hydra_oauth2_consent_request
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
    login_challenge           VARCHAR(40)  NULL REFERENCES hydra_oauth2_authentication_request (challenge) ON DELETE SET NULL,
    requested_at_audience     TEXT         NULL     DEFAULT '',
    acr                       TEXT         NULL     DEFAULT '',
    context                   TEXT         NOT NULL DEFAULT '{}'
);

CREATE INDEX hydra_oauth2_consent_request_client_id_idx ON hydra_oauth2_consent_request (client_id);
CREATE INDEX hydra_oauth2_consent_request_subject_idx ON hydra_oauth2_consent_request (subject);
CREATE INDEX hydra_oauth2_consent_request_login_session_id_idx ON hydra_oauth2_consent_request (login_session_id);
CREATE INDEX hydra_oauth2_consent_request_login_challenge_idx ON hydra_oauth2_consent_request (login_challenge);

CREATE TABLE hydra_oauth2_consent_request_handled
(
    challenge            VARCHAR(40) NOT NULL PRIMARY KEY REFERENCES hydra_oauth2_consent_request (challenge) ON DELETE CASCADE,
    granted_scope        TEXT        NOT NULL,
    remember             INTEGER     NOT NULL,
    remember_for         INTEGER     NOT NULL,
    error                TEXT        NOT NULL,
    requested_at         TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    session_access_token TEXT        NOT NULL,
    session_id_token     TEXT        NOT NULL,
    authenticated_at     TIMESTAMP   NULL,
    was_used             INTEGER     NOT NULL,
    granted_at_audience  TEXT        NULL     DEFAULT '',
    handled_at           TIMESTAMP   NULL
);

CREATE TABLE hydra_oauth2_access
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
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_consent_request_handled (challenge) ON DELETE CASCADE,
    UNIQUE (request_id)
);

CREATE INDEX hydra_oauth2_access_requested_at_idx ON hydra_oauth2_access (requested_at);
CREATE INDEX hydra_oauth2_access_client_id_idx ON hydra_oauth2_access (client_id);
CREATE INDEX hydra_oauth2_access_challenge_id_idx ON hydra_oauth2_access (challenge_id);
CREATE INDEX hydra_oauth2_access_client_id_subject_idx ON hydra_oauth2_access (client_id, subject);

CREATE TABLE hydra_oauth2_authentication_request_handled
(
    challenge                 VARCHAR(40)  NOT NULL PRIMARY KEY REFERENCES hydra_oauth2_authentication_request (challenge) ON DELETE CASCADE,
    subject                   VARCHAR(255) NOT NULL,
    remember                  INTEGER      NOT NULL,
    remember_for              INTEGER      NOT NULL,
    error                     TEXT         NOT NULL,
    acr                       TEXT         NOT NULL,
    requested_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    authenticated_at          TIMESTAMP    NULL,
    was_used                  INTEGER      NOT NULL,
    forced_subject_identifier VARCHAR(255) NULL     DEFAULT '',
    context                   TEXT         NOT NULL DEFAULT '{}'
);

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
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_consent_request_handled (challenge) ON DELETE CASCADE
);

CREATE INDEX hydra_oauth2_code_client_id_idx ON hydra_oauth2_code (client_id);
CREATE INDEX hydra_oauth2_code_challenge_id_idx ON hydra_oauth2_code (challenge_id);

CREATE TABLE hydra_oauth2_jti_blacklist
(
    signature  VARCHAR(64) NOT NULL PRIMARY KEY,
    expires_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX hydra_oauth2_jti_blacklist_expires_at_idx ON hydra_oauth2_jti_blacklist (expires_at);

CREATE TABLE hydra_oauth2_logout_request
(
    challenge    VARCHAR(36)  NOT NULL PRIMARY KEY,
    verifier     VARCHAR(36)  NOT NULL,
    subject      VARCHAR(255) NOT NULL,
    sid          VARCHAR(36)  NOT NULL,
    client_id    VARCHAR(255) NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    request_url  TEXT         NOT NULL,
    redir_url    TEXT         NOT NULL,
    was_used     INTEGER      NOT NULL DEFAULT false,
    accepted     INTEGER      NOT NULL DEFAULT false,
    rejected     INTEGER      NOT NULL DEFAULT false,
    rp_initiated INTEGER      NOT NULL DEFAULT false,
    UNIQUE (verifier)
);

CREATE INDEX hydra_oauth2_logout_request_client_id_idx ON hydra_oauth2_logout_request (client_id);

CREATE TABLE hydra_oauth2_obfuscated_authentication_session
(
    subject            VARCHAR(255) NOT NULL,
    client_id          VARCHAR(255) NOT NULL REFERENCES hydra_client (id) ON DELETE CASCADE,
    subject_obfuscated VARCHAR(255) NOT NULL,
    PRIMARY KEY (subject, client_id)
);

CREATE INDEX hydra_oauth2_obfuscated_authentication_session_client_id_subject_obfuscated_idx ON hydra_oauth2_obfuscated_authentication_session (client_id, subject_obfuscated);

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
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_consent_request_handled (challenge) ON DELETE CASCADE
);

CREATE INDEX hydra_oauth2_oidc_client_id_idx ON hydra_oauth2_oidc (client_id);
CREATE INDEX hydra_oauth2_oidc_challenge_id_idx ON hydra_oauth2_oidc (challenge_id);

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
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_consent_request_handled (challenge) ON DELETE CASCADE
);

CREATE INDEX hydra_oauth2_pkce_client_id_idx ON hydra_oauth2_pkce (client_id);
CREATE INDEX hydra_oauth2_pkce_challenge_id_idx ON hydra_oauth2_pkce (challenge_id);

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
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_consent_request_handled (challenge) ON DELETE CASCADE,
    UNIQUE (request_id)
);

CREATE INDEX hydra_oauth2_refresh_client_id_idx ON hydra_oauth2_refresh (client_id);
CREATE INDEX hydra_oauth2_refresh_challenge_id_idx ON hydra_oauth2_refresh (challenge_id);
CREATE INDEX hydra_oauth2_refresh_client_id_subject_idx ON hydra_oauth2_refresh (client_id, subject);
