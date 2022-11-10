-- hydra_oauth2_jti_blacklist
ALTER TABLE hydra_oauth2_jti_blacklist ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_oauth2_jti_blacklist SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_oauth2_jti_blacklist_tmp" (
    signature  VARCHAR(64) NOT NULL,
    expires_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    nid        CHAR(36)    NOT NULL,
    CHECK (nid != '00000000-0000-0000-0000-000000000000'),
    PRIMARY KEY (signature, nid)
);
INSERT INTO "_hydra_oauth2_jti_blacklist_tmp" (signature, expires_at, nid) SELECT signature, expires_at, nid FROM "hydra_oauth2_jti_blacklist";
DROP TABLE "hydra_oauth2_jti_blacklist";
ALTER TABLE "_hydra_oauth2_jti_blacklist_tmp" RENAME TO "hydra_oauth2_jti_blacklist";

UPDATE hydra_oauth2_jti_blacklist SET nid = (SELECT id FROM networks LIMIT 1);

CREATE INDEX hydra_oauth2_jti_blacklist_expires_at_idx ON hydra_oauth2_jti_blacklist (expires_at, nid);

-- hydra_oauth2_logout_request
ALTER TABLE hydra_oauth2_logout_request ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_oauth2_logout_request SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_oauth2_logout_request_tmp" (
    challenge    VARCHAR(36)  NOT NULL PRIMARY KEY,
    verifier     VARCHAR(36)  NOT NULL,
    subject      VARCHAR(255) NOT NULL,
    sid          VARCHAR(36)  NOT NULL,
    client_id    VARCHAR(255) NULL,
    nid          CHAR(36)     NOT NULL,
    request_url  TEXT         NOT NULL,
    redir_url    TEXT         NOT NULL,
    was_used     INTEGER      NOT NULL DEFAULT false,
    accepted     INTEGER      NOT NULL DEFAULT false,
    rejected     INTEGER      NOT NULL DEFAULT false,
    rp_initiated INTEGER      NOT NULL DEFAULT false,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE,
    UNIQUE (verifier)
);
INSERT INTO "_hydra_oauth2_logout_request_tmp" (challenge, verifier, subject, sid, client_id, request_url, redir_url, was_used, accepted, rejected, rp_initiated, nid) SELECT challenge, verifier, subject, sid, client_id, request_url, redir_url, was_used, accepted, rejected, rp_initiated, nid FROM "hydra_oauth2_logout_request";
DROP TABLE "hydra_oauth2_logout_request";
ALTER TABLE "_hydra_oauth2_logout_request_tmp" RENAME TO "hydra_oauth2_logout_request";

UPDATE hydra_oauth2_logout_request SET nid = (SELECT id FROM networks LIMIT 1);

CREATE INDEX hydra_oauth2_logout_request_client_id_idx ON hydra_oauth2_logout_request (client_id, nid);

-- hydra_oauth2_obfuscated_authentication_session
ALTER TABLE hydra_oauth2_obfuscated_authentication_session ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_oauth2_obfuscated_authentication_session SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_oauth2_obfuscated_authentication_session_tmp" (
    subject            VARCHAR(255) NOT NULL,
    client_id          VARCHAR(255) NOT NULL,
    subject_obfuscated VARCHAR(255) NOT NULL,
    nid                CHAR(36)     NOT NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE,
    PRIMARY KEY (subject, client_id, nid)
);
INSERT INTO "_hydra_oauth2_obfuscated_authentication_session_tmp" (subject, client_id, subject_obfuscated, nid) SELECT subject, client_id, subject_obfuscated, nid FROM "hydra_oauth2_obfuscated_authentication_session";
DROP TABLE "hydra_oauth2_obfuscated_authentication_session";
ALTER TABLE "_hydra_oauth2_obfuscated_authentication_session_tmp" RENAME TO "hydra_oauth2_obfuscated_authentication_session";

UPDATE hydra_oauth2_obfuscated_authentication_session SET nid = (SELECT id FROM networks LIMIT 1);

CREATE UNIQUE INDEX hydra_oauth2_obfuscated_authentication_session_client_id_subject_obfuscated_idx ON hydra_oauth2_obfuscated_authentication_session (client_id, subject_obfuscated, nid);

-- hydra_oauth2_authentication_session
ALTER TABLE hydra_oauth2_authentication_session ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_oauth2_authentication_session SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_oauth2_authentication_session_tmp" (
    id               VARCHAR(40)  NOT NULL PRIMARY KEY,
    authenticated_at TIMESTAMP    NULL,
    subject          VARCHAR(255) NOT NULL,
    nid              CHAR(36)     NOT NULL,
    remember         INTEGER      NOT NULL DEFAULT false,
    CHECK (nid != '00000000-0000-0000-0000-000000000000')
);
INSERT INTO "_hydra_oauth2_authentication_session_tmp" (id, authenticated_at, subject, remember, nid) SELECT id, authenticated_at, subject, remember, nid FROM "hydra_oauth2_authentication_session";
DROP TABLE "hydra_oauth2_authentication_session";
ALTER TABLE "_hydra_oauth2_authentication_session_tmp" RENAME TO "hydra_oauth2_authentication_session";

UPDATE hydra_oauth2_authentication_session SET nid = (SELECT id FROM networks LIMIT 1);

CREATE INDEX hydra_oauth2_authentication_session_subject_idx ON hydra_oauth2_authentication_session (subject, nid);

-- hydra_client
ALTER TABLE hydra_client ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_client SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_client_tmp" (
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
    pk_deprecated                        INTEGER      NULL DEFAULT NULL,
    pk                                   TEXT         PRIMARY KEY,
    audience                             TEXT         NOT NULL,
    created_at                           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    frontchannel_logout_uri              TEXT         NOT NULL DEFAULT '',
    frontchannel_logout_session_required INTEGER      NOT NULL DEFAULT false,
    post_logout_redirect_uris            TEXT         NOT NULL DEFAULT '',
    backchannel_logout_uri               TEXT         NOT NULL DEFAULT '',
    backchannel_logout_session_required  INTEGER      NOT NULL DEFAULT false,
    metadata                             TEXT         NOT NULL DEFAULT '{}',
    token_endpoint_auth_signing_alg      VARCHAR(10)  NOT NULL DEFAULT '',
    registration_access_token_signature  VARCHAR(128) NOT NULL DEFAULT '',
    authorization_code_grant_access_token_lifespan  BIGINT NULL DEFAULT NULL,
    authorization_code_grant_id_token_lifespan      BIGINT NULL DEFAULT NULL,
    authorization_code_grant_refresh_token_lifespan BIGINT NULL DEFAULT NULL,
    client_credentials_grant_access_token_lifespan  BIGINT NULL DEFAULT NULL,
    implicit_grant_access_token_lifespan            BIGINT NULL DEFAULT NULL,
    implicit_grant_id_token_lifespan                BIGINT NULL DEFAULT NULL,
    jwt_bearer_grant_access_token_lifespan          BIGINT NULL DEFAULT NULL,
    password_grant_access_token_lifespan            BIGINT NULL DEFAULT NULL,
    password_grant_refresh_token_lifespan           BIGINT NULL DEFAULT NULL,
    refresh_token_grant_id_token_lifespan           BIGINT NULL DEFAULT NULL,
    refresh_token_grant_access_token_lifespan       BIGINT NULL DEFAULT NULL,
    refresh_token_grant_refresh_token_lifespan      BIGINT NULL DEFAULT NULL,
    nid                                  CHAR(36)     NOT NULL
);
INSERT INTO "_hydra_client_tmp" (
    id,
    client_name,
    client_secret,
    redirect_uris,
    grant_types,
    response_types,
    scope,
    owner,
    policy_uri,
    tos_uri,
    client_uri,
    logo_uri,
    contacts,
    client_secret_expires_at,
    sector_identifier_uri,
    jwks,
    jwks_uri,
    request_uris,
    token_endpoint_auth_method,
    request_object_signing_alg,
    userinfo_signed_response_alg,
    subject_type,
    allowed_cors_origins,
    pk_deprecated,
    pk,
    audience,
    created_at,
    updated_at,
    frontchannel_logout_uri,
    frontchannel_logout_session_required,
    post_logout_redirect_uris,
    backchannel_logout_uri,
    backchannel_logout_session_required,
    metadata,
    token_endpoint_auth_signing_alg,
    registration_access_token_signature,
    authorization_code_grant_access_token_lifespan,
    authorization_code_grant_id_token_lifespan,
    authorization_code_grant_refresh_token_lifespan,
    client_credentials_grant_access_token_lifespan,
    implicit_grant_access_token_lifespan,
    implicit_grant_id_token_lifespan,
    jwt_bearer_grant_access_token_lifespan,
    password_grant_access_token_lifespan,
    password_grant_refresh_token_lifespan,
    refresh_token_grant_id_token_lifespan,
    refresh_token_grant_access_token_lifespan,
    refresh_token_grant_refresh_token_lifespan,
    nid
) SELECT
    id,
    client_name,
    client_secret,
    redirect_uris,
    grant_types,
    response_types,
    scope,
    owner,
    policy_uri,
    tos_uri,
    client_uri,
    logo_uri,
    contacts,
    client_secret_expires_at,
    sector_identifier_uri,
    jwks,
    jwks_uri,
    request_uris,
    token_endpoint_auth_method,
    request_object_signing_alg,
    userinfo_signed_response_alg,
    subject_type,
    allowed_cors_origins,
    pk_deprecated,
    pk,
    audience,
    created_at,
    updated_at,
    frontchannel_logout_uri,
    frontchannel_logout_session_required,
    post_logout_redirect_uris,
    backchannel_logout_uri,
    backchannel_logout_session_required,
    metadata,
    token_endpoint_auth_signing_alg,
    registration_access_token_signature,
    authorization_code_grant_access_token_lifespan,
    authorization_code_grant_id_token_lifespan,
    authorization_code_grant_refresh_token_lifespan,
    client_credentials_grant_access_token_lifespan,
    implicit_grant_access_token_lifespan,
    implicit_grant_id_token_lifespan,
    jwt_bearer_grant_access_token_lifespan,
    password_grant_access_token_lifespan,
    password_grant_refresh_token_lifespan,
    refresh_token_grant_id_token_lifespan,
    refresh_token_grant_access_token_lifespan,
    refresh_token_grant_refresh_token_lifespan,
    nid
FROM "hydra_client";
DROP TABLE "hydra_client";
ALTER TABLE "_hydra_client_tmp" RENAME TO "hydra_client";

UPDATE hydra_client SET nid = (SELECT id FROM networks LIMIT 1);

CREATE UNIQUE INDEX hydra_client_id_nid_uq_idx ON hydra_client (id, nid);
CREATE INDEX hydra_client_id_nid_idx ON hydra_client (id, nid);

-- hydra_oauth2_flow
ALTER TABLE hydra_oauth2_flow ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_oauth2_flow SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_oauth2_flow_tmp" (
    login_challenge           VARCHAR(40)  NOT NULL PRIMARY KEY,
    nid                       CHAR(36)     NOT NULL,
    requested_scope           TEXT         NOT NULL,
    login_verifier            VARCHAR(40)  NOT NULL,
    login_csrf                VARCHAR(40)  NOT NULL,
    subject                   VARCHAR(255) NOT NULL,
    request_url               TEXT         NOT NULL,
    login_skip                INTEGER      NOT NULL,
    client_id                 VARCHAR(255) NOT NULL,
    requested_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    oidc_context              TEXT         NOT NULL,
    login_session_id          VARCHAR(40)  NULL REFERENCES hydra_oauth2_authentication_session (id) ON DELETE SET NULL,
    requested_at_audience     TEXT         NULL DEFAULT '',
    login_initialized_at      TIMESTAMP    NULL,

    state                     INTEGER      NOT NULL,

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
    consent_skip              INTEGER      NULL DEFAULT false,
    consent_verifier          VARCHAR(40)  NULL,
    consent_csrf              VARCHAR(40)  NULL,

    granted_scope             TEXT        NULL,
    granted_at_audience       TEXT        NULL DEFAULT '',
    consent_remember          INTEGER     NULL DEFAULT 0,
    consent_remember_for      INTEGER     NULL,
    consent_handled_at        TIMESTAMP   NULL,
    consent_was_used          INTEGER     NOT NULL DEFAULT false,
    consent_error             TEXT        NULL,
    session_id_token          TEXT        NULL DEFAULT '{}',
    session_access_token      TEXT        NULL DEFAULT '{}',

    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE,

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

INSERT INTO "_hydra_oauth2_flow_tmp" (login_challenge, requested_scope, login_verifier, login_csrf, subject, request_url, login_skip, client_id, requested_at, oidc_context, login_session_id, requested_at_audience, login_initialized_at, state, login_remember, login_remember_for, login_error, acr, login_authenticated_at, login_was_used, forced_subject_identifier, context, amr, consent_challenge_id, consent_skip, consent_verifier, consent_csrf, granted_scope, granted_at_audience, consent_remember, consent_remember_for, consent_handled_at, consent_was_used, consent_error, session_id_token, session_access_token, nid) SELECT login_challenge, requested_scope, login_verifier, login_csrf, subject, request_url, login_skip, client_id, requested_at, oidc_context, login_session_id, requested_at_audience, login_initialized_at, state, login_remember, login_remember_for, login_error, acr, login_authenticated_at, login_was_used, forced_subject_identifier, context, amr, consent_challenge_id, consent_skip, consent_verifier, consent_csrf, granted_scope, granted_at_audience, consent_remember, consent_remember_for, consent_handled_at, consent_was_used, consent_error, session_id_token, session_access_token, nid FROM "hydra_oauth2_flow";
DROP TABLE "hydra_oauth2_flow";
ALTER TABLE "_hydra_oauth2_flow_tmp" RENAME TO "hydra_oauth2_flow";

UPDATE hydra_oauth2_flow SET nid = (SELECT id FROM networks LIMIT 1);

CREATE INDEX hydra_oauth2_flow_client_id_idx ON hydra_oauth2_flow (client_id, nid);
CREATE INDEX hydra_oauth2_flow_login_session_id_idx ON hydra_oauth2_flow (login_session_id);
CREATE INDEX hydra_oauth2_flow_subject_idx ON hydra_oauth2_flow (subject, nid);
CREATE UNIQUE INDEX hydra_oauth2_flow_consent_challenge_id_idx ON hydra_oauth2_flow (consent_challenge_id);
CREATE UNIQUE INDEX hydra_oauth2_flow_login_verifier_idx ON hydra_oauth2_flow (login_verifier);
CREATE UNIQUE INDEX hydra_oauth2_flow_consent_verifier_idx ON hydra_oauth2_flow (consent_verifier);

-- hydra_oauth2_code
ALTER TABLE hydra_oauth2_code ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_oauth2_code SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_oauth2_code_tmp" (
    signature          VARCHAR(255) NOT NULL,
    request_id         VARCHAR(40)  NOT NULL,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          VARCHAR(255) NOT NULL,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL DEFAULT '',
    active             INTEGER      NOT NULL DEFAULT true,
    requested_audience TEXT         NULL     DEFAULT '',
    granted_audience   TEXT         NULL     DEFAULT '',
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE,
    nid                CHAR(36)     NOT NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE
);
INSERT INTO "_hydra_oauth2_code_tmp" (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id, nid) SELECT signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id, nid FROM "hydra_oauth2_code";
DROP TABLE "hydra_oauth2_code";
ALTER TABLE "_hydra_oauth2_code_tmp" RENAME TO "hydra_oauth2_code";

UPDATE hydra_oauth2_code SET nid = (SELECT id FROM networks LIMIT 1);

CREATE INDEX hydra_oauth2_code_client_id_idx ON hydra_oauth2_code (client_id, nid);
CREATE INDEX hydra_oauth2_code_challenge_id_idx ON hydra_oauth2_code (challenge_id, nid);
CREATE INDEX hydra_oauth2_code_request_id_idx ON hydra_oauth2_code (request_id, nid);

-- hydra_oauth2_oidc
ALTER TABLE hydra_oauth2_oidc ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_oauth2_oidc SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_oauth2_oidc_tmp" (
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(40)  NOT NULL,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          VARCHAR(255) NOT NULL,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL DEFAULT '',
    active             INTEGER      NOT NULL DEFAULT true,
    requested_audience TEXT         NULL     DEFAULT '',
    granted_audience   TEXT         NULL     DEFAULT '',
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE,
    nid                CHAR(36)     NOT NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE
);
INSERT INTO "_hydra_oauth2_oidc_tmp" (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id, nid) SELECT signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id, nid FROM "hydra_oauth2_oidc";
DROP TABLE "hydra_oauth2_oidc";
ALTER TABLE "_hydra_oauth2_oidc_tmp" RENAME TO "hydra_oauth2_oidc";

UPDATE hydra_oauth2_oidc SET nid = (SELECT id FROM networks LIMIT 1);

CREATE INDEX hydra_oauth2_oidc_client_id_idx ON hydra_oauth2_oidc (client_id, nid);
CREATE INDEX hydra_oauth2_oidc_challenge_id_idx ON hydra_oauth2_oidc (challenge_id, nid);
CREATE INDEX hydra_oauth2_oidc_request_id_idx ON hydra_oauth2_oidc (request_id, nid);

-- hydra_oauth2_pkce
ALTER TABLE hydra_oauth2_pkce ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_oauth2_pkce SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_oauth2_pkce_tmp" (
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(40)  NOT NULL,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          VARCHAR(255) NOT NULL,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL,
    active             INTEGER      NOT NULL DEFAULT true,
    requested_audience TEXT         NULL     DEFAULT '',
    granted_audience   TEXT         NULL     DEFAULT '',
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE,
    nid                CHAR(36)     NOT NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE
);
INSERT INTO "_hydra_oauth2_pkce_tmp" (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id, nid) SELECT signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id, nid FROM "hydra_oauth2_pkce";
DROP TABLE "hydra_oauth2_pkce";
ALTER TABLE "_hydra_oauth2_pkce_tmp" RENAME TO "hydra_oauth2_pkce";

UPDATE hydra_oauth2_pkce SET nid = (SELECT id FROM networks LIMIT 1);

CREATE INDEX hydra_oauth2_pkce_client_id_idx ON hydra_oauth2_pkce (client_id, nid);
CREATE INDEX hydra_oauth2_pkce_challenge_id_idx ON hydra_oauth2_pkce (challenge_id, nid);
CREATE INDEX hydra_oauth2_pkce_request_id_idx ON hydra_oauth2_pkce (request_id, nid);

-- hydra_oauth2_access
ALTER TABLE hydra_oauth2_access ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_oauth2_access SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_oauth2_access_tmp" (
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(40)  NOT NULL,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          VARCHAR(255) NOT NULL,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL DEFAULT '',
    active             INTEGER      NOT NULL DEFAULT true,
    requested_audience TEXT         NULL     DEFAULT '',
    granted_audience   TEXT         NULL     DEFAULT '',
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE,
    nid                CHAR(36)     NOT NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE
);
INSERT INTO "_hydra_oauth2_access_tmp" (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id, nid) SELECT signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id, nid FROM "hydra_oauth2_access";
DROP TABLE "hydra_oauth2_access";
ALTER TABLE "_hydra_oauth2_access_tmp" RENAME TO "hydra_oauth2_access";

UPDATE hydra_oauth2_access SET nid = (SELECT id FROM networks LIMIT 1);

CREATE INDEX hydra_oauth2_access_requested_at_idx ON hydra_oauth2_access (requested_at, nid);
CREATE INDEX hydra_oauth2_access_client_id_idx ON hydra_oauth2_access (client_id, nid);
CREATE INDEX hydra_oauth2_access_challenge_id_idx ON hydra_oauth2_access (challenge_id, nid);
CREATE INDEX hydra_oauth2_access_client_id_subject_idx ON hydra_oauth2_access (client_id, subject, nid);
CREATE INDEX hydra_oauth2_access_request_id_idx ON hydra_oauth2_access (request_id, nid);

-- hydra_oauth2_refresh
ALTER TABLE hydra_oauth2_refresh ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_oauth2_refresh SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_oauth2_refresh_tmp" (
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(40)  NOT NULL,
    requested_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id          VARCHAR(255) NOT NULL,
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL DEFAULT '',
    active             INTEGER      NOT NULL DEFAULT true,
    requested_audience TEXT         NULL     DEFAULT '',
    granted_audience   TEXT         NULL     DEFAULT '',
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE,
    nid                CHAR(36)     NOT NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE
);
INSERT INTO "_hydra_oauth2_refresh_tmp" (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id, nid) SELECT signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id, nid FROM "hydra_oauth2_refresh";
DROP TABLE "hydra_oauth2_refresh";
ALTER TABLE "_hydra_oauth2_refresh_tmp" RENAME TO "hydra_oauth2_refresh";

UPDATE hydra_oauth2_refresh SET nid = (SELECT id FROM networks LIMIT 1);

CREATE INDEX hydra_oauth2_refresh_client_id_idx ON hydra_oauth2_refresh (client_id, nid);
CREATE INDEX hydra_oauth2_refresh_challenge_id_idx ON hydra_oauth2_refresh (challenge_id, nid);
CREATE INDEX hydra_oauth2_refresh_client_id_subject_idx ON hydra_oauth2_refresh (client_id, subject, nid);
CREATE INDEX hydra_oauth2_refresh_request_id_idx ON hydra_oauth2_refresh (request_id, nid);

-- hydra_oauth2_trusted_jwt_bearer_issuer
ALTER TABLE hydra_oauth2_trusted_jwt_bearer_issuer ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_oauth2_trusted_jwt_bearer_issuer SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_oauth2_trusted_jwt_bearer_issuer_tmp" (
    id         VARCHAR(36) PRIMARY KEY,
    issuer     VARCHAR(255) NOT NULL,
    subject    VARCHAR(255) NOT NULL,
    scope      TEXT         NOT NULL,
    key_set    varchar(255) NOT NULL,
    key_id     varchar(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    nid        CHAR(36)     NOT NULL,
    UNIQUE (issuer, subject, key_id, nid),
    FOREIGN KEY (key_set, key_id, nid) REFERENCES hydra_jwk (sid, kid, nid) ON DELETE CASCADE
);
INSERT INTO "_hydra_oauth2_trusted_jwt_bearer_issuer_tmp" (id, issuer, subject, scope, key_set, key_id, created_at, expires_at, nid) SELECT id, issuer, subject, scope, key_set, key_id, created_at, expires_at, nid FROM "hydra_oauth2_trusted_jwt_bearer_issuer";
DROP TABLE "hydra_oauth2_trusted_jwt_bearer_issuer";
ALTER TABLE "_hydra_oauth2_trusted_jwt_bearer_issuer_tmp" RENAME TO "hydra_oauth2_trusted_jwt_bearer_issuer";

UPDATE hydra_oauth2_trusted_jwt_bearer_issuer SET nid = (SELECT id FROM networks LIMIT 1);

CREATE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_expires_at_idx ON hydra_oauth2_trusted_jwt_bearer_issuer (expires_at);

-- hydra_jwk
ALTER TABLE hydra_jwk ADD COLUMN nid CHAR(36) NULL REFERENCES networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
UPDATE hydra_jwk SET nid = (SELECT id FROM networks LIMIT 1);
CREATE TABLE "_hydra_jwk_tmp" (
    sid             VARCHAR(255) NOT NULL,
    kid             VARCHAR(255) NOT NULL,
    nid             CHAR(36)     NOT NULL,
    version         INTEGER      DEFAULT 0 NOT NULL,
    keydata         TEXT         NOT NULL,
    created_at      TIMESTAMP    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    pk              TEXT         PRIMARY KEY,
    pk_deprecated   INTEGER      NULL DEFAULT NULL,
    CHECK (nid != '00000000-0000-0000-0000-000000000000')
);
INSERT INTO "_hydra_jwk_tmp" (sid, kid, version, keydata, created_at, pk, pk_deprecated, nid) SELECT sid, kid, version, keydata, created_at, pk, pk_deprecated, nid FROM "hydra_jwk";
DROP TABLE "hydra_jwk";
ALTER TABLE "_hydra_jwk_tmp" RENAME TO "hydra_jwk";

UPDATE hydra_jwk SET nid = (SELECT id FROM networks LIMIT 1);

CREATE UNIQUE INDEX hydra_jwk_sid_kid_nid_key ON hydra_jwk (sid, kid, nid);
