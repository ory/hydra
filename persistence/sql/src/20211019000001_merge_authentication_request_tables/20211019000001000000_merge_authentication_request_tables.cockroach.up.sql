CREATE TABLE hydra_oauth2_flow
(
    login_challenge           character varying(40) NOT NULL,
    requested_scope           text NOT NULL DEFAULT '[]',
    login_verifier            character varying(40) NOT NULL,
    login_csrf                character varying(40) NOT NULL,
    subject                   character varying(255) NOT NULL,
    request_url               text NOT NULL,
    login_skip                boolean NOT NULL,
    client_id                 character varying(255) NOT NULL,
    requested_at              timestamp without time zone DEFAULT now() NOT NULL,
    login_initialized_at      timestamp without time zone NULL DEFAULT NULL,
    oidc_context              jsonb NOT NULL DEFAULT '{}',
    login_session_id          character varying(40) NULL,
    requested_at_audience     text NULL DEFAULT '[]',

    state                     INTEGER      NOT NULL,

    login_remember boolean NOT NULL DEFAULT false,
    login_remember_for integer NOT NULL,
    login_error text NULL,
    acr text  NOT NULL DEFAULT '',
    login_authenticated_at timestamp without time zone NULL DEFAULT NULL,
    login_was_used boolean NOT NULL DEFAULT false,
    forced_subject_identifier character varying(255) NOT NULL DEFAULT ''::character varying,
    context jsonb DEFAULT '{}',
    amr text DEFAULT '[]',

    consent_challenge_id character varying(40) NULL,
    consent_skip boolean DEFAULT false NOT NULL,
    consent_verifier character varying(40) NULL,
    consent_csrf character varying(40) NULL,

    granted_scope text NOT NULL DEFAULT '[]',
    granted_at_audience text NOT NULL DEFAULT '[]',
    consent_remember boolean DEFAULT false NOT NULL,
    consent_remember_for integer NULL,
    consent_handled_at TIMESTAMP WITHOUT TIME ZONE NULL,
    consent_error TEXT NULL,
    session_access_token jsonb DEFAULT '{}' NOT NULL,
    session_id_token jsonb DEFAULT '{}' NOT NULL,
    consent_was_used boolean DEFAULT false NOT NULL,

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
    cast(coalesce(hydra_oauth2_authentication_request.oidc_context, '{}') as jsonb),
    hydra_oauth2_authentication_request.login_session_id,
    hydra_oauth2_authentication_request.requested_at_audience,

    coalesce(hydra_oauth2_authentication_request_handled.remember, false),
    coalesce(hydra_oauth2_authentication_request_handled.remember_for, 0),
    hydra_oauth2_authentication_request_handled.error,
    coalesce(hydra_oauth2_authentication_request_handled.acr, ''),
    hydra_oauth2_authentication_request_handled.authenticated_at,
    coalesce(hydra_oauth2_authentication_request_handled.was_used, false),
    coalesce(hydra_oauth2_consent_request.forced_subject_identifier, hydra_oauth2_authentication_request_handled.forced_subject_identifier, ''),
    cast(coalesce(hydra_oauth2_authentication_request_handled.context, '{}') as jsonb),
    coalesce(hydra_oauth2_authentication_request_handled.amr, ''),

    hydra_oauth2_consent_request.challenge,
    hydra_oauth2_consent_request.verifier,
    coalesce(hydra_oauth2_consent_request.skip, false),
    hydra_oauth2_consent_request.csrf,

    coalesce(hydra_oauth2_consent_request_handled.granted_scope, '[]'),
    coalesce(hydra_oauth2_consent_request_handled.remember, false),
    coalesce(hydra_oauth2_consent_request_handled.remember_for, 0),
    hydra_oauth2_consent_request_handled.error,
    cast(coalesce(hydra_oauth2_consent_request_handled.session_access_token, '{}') as jsonb),
    cast(coalesce(hydra_oauth2_consent_request_handled.session_id_token, '{}') as jsonb),
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
CREATE INDEX hydra_oauth2_flow_client_id_subject_idx ON hydra_oauth2_flow USING btree (client_id, subject);
CREATE INDEX hydra_oauth2_flow_cid_idx ON hydra_oauth2_flow USING btree (client_id);
CREATE INDEX hydra_oauth2_flow_login_session_id_idx ON hydra_oauth2_flow USING btree (login_session_id);
CREATE INDEX hydra_oauth2_flow_sub_idx ON hydra_oauth2_flow USING btree (subject);
CREATE UNIQUE INDEX hydra_oauth2_flow_consent_challenge_idx ON hydra_oauth2_flow USING btree (consent_challenge_id);
CREATE UNIQUE INDEX hydra_oauth2_flow_login_verifier_idx ON hydra_oauth2_flow USING btree (login_verifier);
CREATE INDEX hydra_oauth2_flow_consent_verifier_idx ON hydra_oauth2_flow USING btree (consent_verifier);
--split
ALTER TABLE ONLY hydra_oauth2_flow ADD CONSTRAINT hydra_oauth2_flow_pkey PRIMARY KEY (login_challenge);
--split
ALTER TABLE ONLY hydra_oauth2_flow ADD CONSTRAINT hydra_oauth2_flow_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
ALTER TABLE ONLY hydra_oauth2_flow ADD CONSTRAINT hydra_oauth2_flow_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES hydra_oauth2_authentication_session(id) ON DELETE CASCADE;

ALTER TABLE ONLY hydra_oauth2_access DROP CONSTRAINT hydra_oauth2_access_challenge_id_fk;
ALTER TABLE ONLY hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;

ALTER TABLE ONLY hydra_oauth2_code DROP CONSTRAINT hydra_oauth2_code_challenge_id_fk;
ALTER TABLE ONLY hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;

ALTER TABLE ONLY hydra_oauth2_oidc DROP CONSTRAINT hydra_oauth2_oidc_challenge_id_fk;
ALTER TABLE ONLY hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;

ALTER TABLE ONLY hydra_oauth2_pkce DROP CONSTRAINT hydra_oauth2_pkce_challenge_id_fk;
ALTER TABLE ONLY hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;

ALTER TABLE ONLY hydra_oauth2_refresh DROP CONSTRAINT hydra_oauth2_refresh_challenge_id_fk;
ALTER TABLE ONLY hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;

--split
DROP TABLE hydra_oauth2_consent_request_handled;
--split
DROP TABLE hydra_oauth2_consent_request;
--split
DROP TABLE hydra_oauth2_authentication_request_handled;
--split
DROP TABLE hydra_oauth2_authentication_request;
