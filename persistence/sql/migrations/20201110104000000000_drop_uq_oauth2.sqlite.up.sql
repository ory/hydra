CREATE TABLE "_hydra_oauth2_access_tmp"
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

INSERT INTO "_hydra_oauth2_access_tmp" (
    signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id
) SELECT signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id FROM "hydra_oauth2_access";

DROP TABLE "hydra_oauth2_access";
ALTER TABLE "_hydra_oauth2_access_tmp" RENAME TO "hydra_oauth2_access";

CREATE INDEX hydra_oauth2_access_requested_at_idx ON hydra_oauth2_access (requested_at);
CREATE INDEX hydra_oauth2_access_client_id_idx ON hydra_oauth2_access (client_id);
CREATE INDEX hydra_oauth2_access_challenge_id_idx ON hydra_oauth2_access (challenge_id);
CREATE INDEX hydra_oauth2_access_client_id_subject_idx ON hydra_oauth2_access (client_id, subject);

---

CREATE TABLE "_hydra_oauth2_refresh_tmp"
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

INSERT INTO "_hydra_oauth2_refresh_tmp" (
    signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id
) SELECT signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id FROM "hydra_oauth2_refresh";

DROP TABLE "hydra_oauth2_refresh";
ALTER TABLE "_hydra_oauth2_refresh_tmp" RENAME TO "hydra_oauth2_refresh";

CREATE INDEX hydra_oauth2_refresh_client_id_idx ON hydra_oauth2_refresh (client_id);
CREATE INDEX hydra_oauth2_refresh_challenge_id_idx ON hydra_oauth2_refresh (challenge_id);
CREATE INDEX hydra_oauth2_refresh_client_id_subject_idx ON hydra_oauth2_refresh (client_id, subject);

---

CREATE INDEX hydra_oauth2_access_request_id_idx ON hydra_oauth2_access (request_id);
CREATE INDEX hydra_oauth2_refresh_request_id_idx ON hydra_oauth2_refresh (request_id);

CREATE INDEX hydra_oauth2_code_request_id_idx ON hydra_oauth2_code (request_id);
CREATE INDEX hydra_oauth2_oidc_request_id_idx ON hydra_oauth2_oidc (request_id);
CREATE INDEX hydra_oauth2_pkce_request_id_idx ON hydra_oauth2_pkce (request_id);
