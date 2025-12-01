-- migrations hash: f5af1bbf8624fd09cf094c1a5745a255e6ea8d56bc7177e0db9eb73d72f1b3dae0fdf3312458a84550c1f2374b0a1ea6fdd026e22267fbf1808b8e8611fb39c0

CREATE TABLE "hydra_client"
(
  id                                              VARCHAR(255) NOT NULL,
  client_name                                     TEXT         NOT NULL,
  client_secret                                   TEXT         NOT NULL,
  redirect_uris                                   TEXT         NOT NULL,
  grant_types                                     TEXT         NOT NULL,
  response_types                                  TEXT         NOT NULL,
  scope                                           TEXT         NOT NULL,
  owner                                           TEXT         NOT NULL,
  policy_uri                                      TEXT         NOT NULL,
  tos_uri                                         TEXT         NOT NULL,
  client_uri                                      TEXT         NOT NULL,
  logo_uri                                        TEXT         NOT NULL,
  contacts                                        TEXT         NOT NULL,
  client_secret_expires_at                        INTEGER      NOT NULL DEFAULT 0,
  sector_identifier_uri                           TEXT         NOT NULL,
  jwks                                            TEXT         NOT NULL,
  jwks_uri                                        TEXT         NOT NULL,
  request_uris                                    TEXT         NOT NULL,
  token_endpoint_auth_method                      VARCHAR(25)  NOT NULL DEFAULT '',
  request_object_signing_alg                      VARCHAR(10)  NOT NULL DEFAULT '',
  userinfo_signed_response_alg                    VARCHAR(10)  NOT NULL DEFAULT '',
  subject_type                                    VARCHAR(15)  NOT NULL DEFAULT '',
  allowed_cors_origins                            TEXT         NOT NULL,
  pk                                              TEXT         NULL,
  pk_deprecated                                   INTEGER NULL DEFAULT NULL,
  audience                                        TEXT         NOT NULL,
  created_at                                      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at                                      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  frontchannel_logout_uri                         TEXT         NOT NULL DEFAULT '',
  frontchannel_logout_session_required            INTEGER      NOT NULL DEFAULT false,
  post_logout_redirect_uris                       TEXT         NOT NULL DEFAULT '',
  backchannel_logout_uri                          TEXT         NOT NULL DEFAULT '',
  backchannel_logout_session_required             INTEGER      NOT NULL DEFAULT false,
  metadata                                        TEXT         NOT NULL DEFAULT '{}',
  token_endpoint_auth_signing_alg                 VARCHAR(10)  NOT NULL DEFAULT '',
  registration_access_token_signature             VARCHAR(128) NOT NULL DEFAULT '',
  access_token_strategy                           VARCHAR(10)  NOT NULL DEFAULT '',
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
  skip_consent                                    BOOLEAN      NOT NULL DEFAULT false,
  nid                                             CHAR(36)     NOT NULL, skip_logout_consent BOOLEAN NULL, device_authorization_grant_id_token_lifespan BIGINT NULL DEFAULT NULL, device_authorization_grant_access_token_lifespan BIGINT NULL DEFAULT NULL, device_authorization_grant_refresh_token_lifespan BIGINT NULL DEFAULT NULL,
  PRIMARY KEY (id, nid)
);
CREATE TABLE "hydra_jwk" (
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
CREATE INDEX hydra_jwk_nid_sid_created_at_idx ON hydra_jwk (nid, sid, created_at);
CREATE INDEX hydra_jwk_nid_sid_kid_created_at_idx ON hydra_jwk (nid, sid, kid, created_at);
CREATE UNIQUE INDEX hydra_jwk_sid_kid_nid_key ON hydra_jwk (sid, kid, nid);
CREATE TABLE "hydra_oauth2_access" (
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
    nid                CHAR(36)     NOT NULL, expires_at TIMESTAMP NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE
);
CREATE INDEX hydra_oauth2_access_challenge_id_idx ON hydra_oauth2_access (challenge_id, nid);
CREATE INDEX hydra_oauth2_access_client_id_idx ON hydra_oauth2_access (client_id, nid);
CREATE INDEX hydra_oauth2_access_request_id_idx ON hydra_oauth2_access (request_id, nid);
CREATE INDEX hydra_oauth2_access_requested_at_idx ON hydra_oauth2_access (requested_at, nid);
CREATE TABLE "hydra_oauth2_authentication_session" (
    id               VARCHAR(40)  NOT NULL PRIMARY KEY,
    authenticated_at TIMESTAMP    NULL,
    subject          VARCHAR(255) NOT NULL,
    nid              CHAR(36)     NOT NULL,
    remember         INTEGER      NOT NULL DEFAULT false, identity_provider_session_id VARCHAR(40), expires_at TIMESTAMP NULL,
    CHECK (nid != '00000000-0000-0000-0000-000000000000')
);
CREATE INDEX hydra_oauth2_authentication_session_subject_idx ON hydra_oauth2_authentication_session (subject, nid);
CREATE TABLE "hydra_oauth2_code" (
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
    nid                CHAR(36)     NOT NULL, expires_at TIMESTAMP NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE
);
CREATE INDEX hydra_oauth2_code_challenge_id_idx ON hydra_oauth2_code (challenge_id, nid);
CREATE INDEX hydra_oauth2_code_client_id_idx ON hydra_oauth2_code (client_id, nid);
CREATE TABLE hydra_oauth2_device_auth_codes
(
  device_code_signature VARCHAR(255)  NOT NULL,
  user_code_signature   VARCHAR(255)  NOT NULL,
  request_id            VARCHAR(40)   NOT NULL,
  requested_at          TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  client_id             VARCHAR(255)  NOT NULL,
  scope                 VARCHAR(1024) NOT NULL,
  granted_scope         VARCHAR(1024) NOT NULL,
  form_data             VARCHAR(4096) NOT NULL,
  session_data          TEXT NOT NULL,
  subject               VARCHAR(255)  NOT NULL DEFAULT '',
  device_code_active    BOOL          NOT NULL DEFAULT true,
  user_code_state       SMALLINT      NOT NULL DEFAULT 0,
  requested_audience    VARCHAR(1024) NOT NULL,
  granted_audience      VARCHAR(1024) NOT NULL,
  challenge_id          VARCHAR(40)   NULL,
  expires_at            TIMESTAMP     NULL,
  nid                   UUID          NOT NULL,

  FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE,
  FOREIGN KEY (nid) REFERENCES networks (id) ON UPDATE RESTRICT ON DELETE CASCADE,
  FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow (consent_challenge_id) ON DELETE CASCADE,
  PRIMARY KEY (device_code_signature, nid)
);
CREATE INDEX hydra_oauth2_device_auth_codes_challenge_id_idx ON hydra_oauth2_device_auth_codes (challenge_id);
CREATE INDEX hydra_oauth2_device_auth_codes_client_id_idx ON hydra_oauth2_device_auth_codes (client_id, nid);
CREATE INDEX hydra_oauth2_device_auth_codes_request_id_idx ON hydra_oauth2_device_auth_codes (request_id, nid);
CREATE UNIQUE INDEX hydra_oauth2_device_auth_codes_user_code_signature_idx ON hydra_oauth2_device_auth_codes (nid, user_code_signature);
CREATE TABLE "hydra_oauth2_flow" (
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
    session_access_token      TEXT        NULL DEFAULT '{}', login_extend_session_lifespan BOOLEAN NOT NULL DEFAULT FALSE, identity_provider_session_id VARCHAR(40), device_challenge_id VARCHAR(255) NULL, device_code_request_id VARCHAR(255) NULL, device_verifier VARCHAR(40) NULL, device_csrf VARCHAR(40) NULL, device_was_used BOOLEAN NULL, device_handled_at TIMESTAMP NULL, device_error VARCHAR(2048) NULL, expires_at TIMESTAMP
    GENERATED ALWAYS AS (if(consent_remember_for > 0, datetime(requested_at, '+' || consent_remember_for || ' seconds'), NULL)) VIRTUAL,

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
CREATE INDEX hydra_oauth2_flow_client_id_idx ON hydra_oauth2_flow (client_id, nid);
CREATE UNIQUE INDEX hydra_oauth2_flow_consent_challenge_id_idx ON hydra_oauth2_flow (consent_challenge_id);
CREATE UNIQUE INDEX hydra_oauth2_flow_device_challenge_idx ON hydra_oauth2_flow (device_challenge_id);
CREATE INDEX hydra_oauth2_flow_login_session_id_idx ON hydra_oauth2_flow (login_session_id);
CREATE INDEX hydra_oauth2_flow_previous_consents_idx
  ON hydra_oauth2_flow (subject, client_id, nid, consent_skip, consent_error, consent_remember);
CREATE INDEX hydra_oauth2_flow_subject_idx ON hydra_oauth2_flow (subject, nid);
CREATE TABLE "hydra_oauth2_jti_blacklist" (
    signature  VARCHAR(64) NOT NULL,
    expires_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    nid        CHAR(36)    NOT NULL,
    CHECK (nid != '00000000-0000-0000-0000-000000000000'),
    PRIMARY KEY (signature, nid)
);
CREATE INDEX hydra_oauth2_jti_blacklist_expires_at_idx ON hydra_oauth2_jti_blacklist (expires_at, nid);
CREATE TABLE "hydra_oauth2_logout_request" (
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
    rp_initiated INTEGER      NOT NULL DEFAULT false, expires_at timestamp NULL, requested_at timestamp NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE,
    UNIQUE (verifier)
);
CREATE INDEX hydra_oauth2_logout_request_client_id_idx ON hydra_oauth2_logout_request (client_id, nid);
CREATE TABLE "hydra_oauth2_obfuscated_authentication_session" (
    subject            VARCHAR(255) NOT NULL,
    client_id          VARCHAR(255) NOT NULL,
    subject_obfuscated VARCHAR(255) NOT NULL,
    nid                CHAR(36)     NOT NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE,
    PRIMARY KEY (subject, client_id, nid)
);
CREATE UNIQUE INDEX hydra_oauth2_obfuscated_authentication_session_client_id_subject_obfuscated_idx ON hydra_oauth2_obfuscated_authentication_session (client_id, subject_obfuscated, nid);
CREATE TABLE "hydra_oauth2_oidc" (
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
    nid                CHAR(36)     NOT NULL, expires_at TIMESTAMP NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE
);
CREATE INDEX hydra_oauth2_oidc_challenge_id_idx ON hydra_oauth2_oidc (challenge_id, nid);
CREATE INDEX hydra_oauth2_oidc_client_id_idx ON hydra_oauth2_oidc (client_id, nid);
CREATE TABLE "hydra_oauth2_pkce" (
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
    nid                CHAR(36)     NOT NULL, expires_at TIMESTAMP NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE
);
CREATE INDEX hydra_oauth2_pkce_challenge_id_idx ON hydra_oauth2_pkce (challenge_id, nid);
CREATE INDEX hydra_oauth2_pkce_client_id_idx ON hydra_oauth2_pkce (client_id, nid);
CREATE TABLE "hydra_oauth2_refresh" (
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
    nid                CHAR(36)     NOT NULL, expires_at TIMESTAMP NULL, first_used_at TIMESTAMP DEFAULT NULL, access_token_signature VARCHAR(255) DEFAULT NULL, used_times INT NULL,
    FOREIGN KEY (client_id, nid) REFERENCES hydra_client (id, nid) ON DELETE CASCADE
);
CREATE INDEX hydra_oauth2_refresh_challenge_id_idx ON hydra_oauth2_refresh (challenge_id, nid);
CREATE INDEX hydra_oauth2_refresh_client_id_idx ON hydra_oauth2_refresh (client_id, nid);
CREATE INDEX hydra_oauth2_refresh_request_id_idx ON hydra_oauth2_refresh (request_id, nid);
CREATE INDEX hydra_oauth2_refresh_requested_at_idx ON hydra_oauth2_refresh (nid, requested_at);
CREATE TABLE "hydra_oauth2_trusted_jwt_bearer_issuer" (
    id         VARCHAR(36) PRIMARY KEY,
    issuer     VARCHAR(255) NOT NULL,
    subject    VARCHAR(255) NOT NULL,
    scope      TEXT         NOT NULL,
    key_set    varchar(255) NOT NULL,
    key_id     varchar(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    nid        CHAR(36)     NOT NULL, allow_any_subject INTEGER NOT NULL DEFAULT FALSE,
    UNIQUE (issuer, subject, key_id, nid),
    FOREIGN KEY (key_set, key_id, nid) REFERENCES hydra_jwk (sid, kid, nid) ON DELETE CASCADE
);
CREATE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_expires_at_idx ON hydra_oauth2_trusted_jwt_bearer_issuer (expires_at);
CREATE UNIQUE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_nid_uq_idx ON hydra_oauth2_trusted_jwt_bearer_issuer (nid ASC, key_id ASC, issuer ASC, subject ASC);
CREATE TABLE "networks" (
  "id" TEXT PRIMARY KEY,
  "created_at" DATETIME NOT NULL,
  "updated_at" DATETIME NOT NULL
);
CREATE TABLE schema_migration (version VARCHAR (48) NOT NULL, version_self INT NOT NULL DEFAULT 0);
CREATE UNIQUE INDEX schema_migration_version_idx ON schema_migration (version);
CREATE INDEX schema_migration_version_self_idx ON schema_migration (version_self);
