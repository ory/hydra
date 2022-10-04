CREATE TABLE hydra_oauth2_flow
(
    `login_challenge` varchar(40) NOT NULL,
    `requested_scope` text NOT NULL DEFAULT ('[]'),
    `login_verifier` varchar(40) NOT NULL,
    `login_csrf` varchar(40) NOT NULL,
    `subject` varchar(255) NOT NULL,
    `request_url` text NOT NULL,
    `login_skip` tinyint(1) NOT NULL,
    `client_id` varchar(255) NOT NULL,
    `requested_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `login_initialized_at` timestamp NULL DEFAULT NULL,
    `oidc_context` json NOT NULL DEFAULT (('{}')),
    `login_session_id` varchar(40) NULL,
    `requested_at_audience` text NULL DEFAULT ('[]'),

    `state` smallint NOT NULL,

    `login_remember` tinyint(1) NOT NULL DEFAULT false,
    `login_remember_for` int(11) NOT NULL,
    `login_error` text NULL,
    `acr` text  NOT NULL DEFAULT (''),
    `login_authenticated_at` timestamp NULL DEFAULT NULL,
    `login_was_used` tinyint(1) NOT NULL DEFAULT false,
    `forced_subject_identifier` varchar(255) NOT NULL DEFAULT '',
    `context` json NOT NULL DEFAULT ('{}'),
    `amr` text NOT NULL DEFAULT ('[]'),

    `consent_challenge_id` varchar(40) NULL,
    `consent_skip` tinyint(1) NOT NULL DEFAULT 0,
    `consent_verifier` varchar(40) NULL,
    `consent_csrf` varchar(40) NULL,

    `granted_scope` text NOT NULL DEFAULT ('[]'),
    `granted_at_audience` text NOT NULL DEFAULT ('[]'),
    `consent_remember` tinyint(1) NOT NULL DEFAULT false,
    `consent_remember_for` int(11) NULL,
    `consent_handled_at` timestamp NULL DEFAULT NULL,
    `consent_error` TEXT NULL,
    `session_access_token` json DEFAULT ('{}') NOT NULL,
    `session_id_token` json DEFAULT ('{}') NOT NULL,
    `consent_was_used` tinyint(1),

    PRIMARY KEY (`login_challenge`),
    UNIQUE KEY `hydra_oauth2_flow_login_verifier_idx` (`login_verifier`),
    KEY `hydra_oauth2_flow_cid_idx` (`client_id`),
    KEY `hydra_oauth2_flow_sub_idx` (`subject`),
    KEY `hydra_oauth2_flow_login_session_id_idx` (`login_session_id`),
    CONSTRAINT `hydra_oauth2_flow_client_id_fk` FOREIGN KEY (`client_id`) REFERENCES `hydra_client` (`id`) ON DELETE CASCADE,
    CONSTRAINT `hydra_oauth2_flow_login_session_id_fk` FOREIGN KEY (`login_session_id`) REFERENCES `hydra_oauth2_authentication_session` (`id`) ON DELETE CASCADE,

    UNIQUE KEY `hydra_oauth2_flow_consent_challenge_idx` (`consent_challenge_id`),
    KEY `hydra_oauth2_flow_consent_verifier_idx` (`consent_verifier`),
    KEY `hydra_oauth2_flow_client_id_subject_idx` (`client_id`,`subject`)
);

ALTER TABLE hydra_oauth2_flow ADD CONSTRAINT hydra_oauth2_flow_chk CHECK (
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

    `consent_challenge_id`,
    `consent_verifier`,
    `consent_skip`,
    `consent_csrf`,

    `granted_scope`,
    `consent_remember`,
    `consent_remember_for`,
    `consent_error`,
    `session_access_token`,
    `session_id_token`,
    `consent_was_used`,
    `granted_at_audience`,
    `consent_handled_at`
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
    coalesce(hydra_oauth2_consent_request_handled.session_access_token, ('{}')),
    coalesce(hydra_oauth2_consent_request_handled.session_id_token, ('{}')),
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
ALTER TABLE hydra_oauth2_access DROP FOREIGN KEY hydra_oauth2_access_challenge_id_fk;
--split
ALTER TABLE hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;

--split
ALTER TABLE hydra_oauth2_code DROP FOREIGN KEY hydra_oauth2_code_challenge_id_fk;
--split
ALTER TABLE hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;

--split
ALTER TABLE hydra_oauth2_oidc DROP FOREIGN KEY hydra_oauth2_oidc_challenge_id_fk;
--split
ALTER TABLE hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;

--split
ALTER TABLE hydra_oauth2_pkce DROP FOREIGN KEY hydra_oauth2_pkce_challenge_id_fk;
--split
ALTER TABLE hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;

--split
ALTER TABLE hydra_oauth2_refresh DROP FOREIGN KEY hydra_oauth2_refresh_challenge_id_fk;
--split
ALTER TABLE hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;

--split
DROP TABLE hydra_oauth2_consent_request_handled;
--split
DROP TABLE hydra_oauth2_consent_request;
--split
DROP TABLE hydra_oauth2_authentication_request_handled;
--split
DROP TABLE hydra_oauth2_authentication_request;
