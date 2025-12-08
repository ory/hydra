CREATE TABLE hydra_oauth2_flow_prev (
  login_challenge               VARCHAR(40)   NOT NULL PRIMARY KEY,
  nid                           CHAR(36)      NOT NULL,
  requested_scope               TEXT          NOT NULL,
  login_verifier                VARCHAR(40)   NOT NULL,
  login_csrf                    VARCHAR(40)   NOT NULL,
  subject                       VARCHAR(255)  NOT NULL,
  request_url                   TEXT          NOT NULL,
  login_skip                    INTEGER       NOT NULL,
  client_id                     VARCHAR(255)  NOT NULL,
  requested_at                  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  oidc_context                  TEXT          NOT NULL,
  login_session_id              VARCHAR(40)   NULL REFERENCES hydra_oauth2_authentication_session (id) ON DELETE SET NULL,
  requested_at_audience         TEXT          NULL,
  login_initialized_at          TIMESTAMP     NULL,

  state                         INTEGER       NOT NULL,

  login_remember                INTEGER       NULL,
  login_remember_for            INTEGER       NULL,
  login_error                   TEXT          NULL,
  acr                           TEXT          NULL,
  login_authenticated_at        TIMESTAMP     NULL,
  login_was_used                INTEGER       NULL,
  forced_subject_identifier     VARCHAR(255)  NULL,
  context                       TEXT          NULL,
  amr                           TEXT          NULL,

  consent_challenge_id          VARCHAR(40)   NULL,
  consent_skip                  INTEGER       NULL,
  consent_verifier              VARCHAR(40)   NULL,
  consent_csrf                  VARCHAR(40)   NULL,

  granted_scope                 TEXT          NULL,
  granted_at_audience           TEXT          NULL,
  consent_remember              INTEGER       NULL,
  consent_remember_for          INTEGER       NULL,
  consent_handled_at            TIMESTAMP     NULL,
  consent_was_used              INTEGER       NOT NULL,
  consent_error                 TEXT          NULL,
  session_id_token              TEXT          NULL,
  session_access_token          TEXT          NULL,
  login_extend_session_lifespan BOOLEAN       NULL,
  identity_provider_session_id  VARCHAR(40)   NULL,
  device_challenge_id           VARCHAR(255)  NULL,
  device_code_request_id        VARCHAR(255)  NULL,
  device_verifier               VARCHAR(40)   NULL,
  device_csrf                   VARCHAR(40)   NULL,
  device_was_used               BOOLEAN       NULL,
  device_handled_at             TIMESTAMP     NULL,
  device_error                  VARCHAR(2048) NULL,
  expires_at                    TIMESTAMP GENERATED ALWAYS AS (IF(consent_remember_for > 0,
                                                                  datetime(requested_at, '+' || consent_remember_for || ' seconds'),
                                                                  NULL)) VIRTUAL,

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

INSERT
INTO hydra_oauth2_flow_prev (login_challenge, nid, requested_scope, login_verifier, login_csrf, subject, request_url,
                             login_skip, client_id, requested_at, oidc_context, login_session_id, requested_at_audience,
                             login_initialized_at, state, login_remember, login_remember_for, login_error, acr,
                             login_authenticated_at, login_was_used, forced_subject_identifier, context, amr,
                             consent_challenge_id, consent_skip, consent_verifier, consent_csrf, granted_scope,
                             granted_at_audience, consent_remember, consent_remember_for, consent_handled_at,
                             consent_was_used, consent_error, session_id_token, session_access_token,
                             login_extend_session_lifespan, identity_provider_session_id, device_challenge_id,
                             device_code_request_id, device_verifier, device_csrf, device_was_used, device_handled_at,
                             device_error)
SELECT login_challenge,
       nid,
       requested_scope,
       login_verifier,
       login_csrf,
       subject,
       request_url,
       login_skip,
       client_id,
       requested_at,
       oidc_context,
       login_session_id,
       requested_at_audience,
       login_initialized_at,
       state,
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
       consent_skip,
       consent_verifier,
       consent_csrf,
       granted_scope,
       granted_at_audience,
       consent_remember,
       consent_remember_for,
       consent_handled_at,
       consent_was_used,
       consent_error,
       session_id_token,
       session_access_token,
       login_extend_session_lifespan,
       identity_provider_session_id,
       device_challenge_id,
       device_code_request_id,
       device_verifier,
       device_csrf,
       device_was_used,
       device_handled_at,
       device_error
FROM hydra_oauth2_flow;

DROP TABLE hydra_oauth2_flow;

ALTER TABLE hydra_oauth2_flow_prev
  RENAME TO hydra_oauth2_flow;

CREATE INDEX hydra_oauth2_flow_client_id_idx ON hydra_oauth2_flow (client_id, nid);
CREATE INDEX hydra_oauth2_flow_login_session_id_idx ON hydra_oauth2_flow (login_session_id);
CREATE INDEX hydra_oauth2_flow_subject_idx ON hydra_oauth2_flow (subject, nid);
CREATE UNIQUE INDEX hydra_oauth2_flow_consent_challenge_id_idx ON hydra_oauth2_flow (consent_challenge_id);
CREATE INDEX hydra_oauth2_flow_previous_consents_idx ON hydra_oauth2_flow (subject, client_id, nid, consent_skip,
                                                                           consent_error, consent_remember);
CREATE UNIQUE INDEX hydra_oauth2_flow_device_challenge_idx ON hydra_oauth2_flow (device_challenge_id);
