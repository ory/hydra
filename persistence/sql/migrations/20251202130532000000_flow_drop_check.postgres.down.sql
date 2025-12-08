ALTER TABLE hydra_oauth2_flow
  ADD CONSTRAINT hydra_oauth2_flow_check CHECK (((state = 128) OR (state = 129) OR (state = 1) OR ((state = 2) AND
                                                                                                   ((login_remember IS NOT NULL) AND
                                                                                                    (login_remember_for IS NOT NULL) AND
                                                                                                    (login_error IS NOT NULL) AND
                                                                                                    (acr IS NOT NULL) AND
                                                                                                    (login_was_used IS NOT NULL) AND
                                                                                                    (context IS NOT NULL) AND
                                                                                                    (amr IS NOT NULL))) OR
                                                 ((state = 3) AND
                                                  ((login_remember IS NOT NULL) AND (login_remember_for IS NOT NULL) AND
                                                   (login_error IS NOT NULL) AND (acr IS NOT NULL) AND
                                                   (login_was_used IS NOT NULL) AND (context IS NOT NULL) AND
                                                   (amr IS NOT NULL))) OR ((state = 4) AND
                                                                           ((login_remember IS NOT NULL) AND
                                                                            (login_remember_for IS NOT NULL) AND
                                                                            (login_error IS NOT NULL) AND
                                                                            (acr IS NOT NULL) AND
                                                                            (login_was_used IS NOT NULL) AND
                                                                            (context IS NOT NULL) AND
                                                                            (amr IS NOT NULL) AND
                                                                            (consent_challenge_id IS NOT NULL) AND
                                                                            (consent_verifier IS NOT NULL) AND
                                                                            (consent_skip IS NOT NULL) AND
                                                                            (consent_csrf IS NOT NULL))) OR
                                                 ((state = 5) AND
                                                  ((login_remember IS NOT NULL) AND (login_remember_for IS NOT NULL) AND
                                                   (login_error IS NOT NULL) AND (acr IS NOT NULL) AND
                                                   (login_was_used IS NOT NULL) AND (context IS NOT NULL) AND
                                                   (amr IS NOT NULL) AND (consent_challenge_id IS NOT NULL) AND
                                                   (consent_verifier IS NOT NULL) AND (consent_skip IS NOT NULL) AND
                                                   (consent_csrf IS NOT NULL))) OR ((state = 6) AND
                                                                                    ((login_remember IS NOT NULL) AND
                                                                                     (login_remember_for IS NOT NULL) AND
                                                                                     (login_error IS NOT NULL) AND
                                                                                     (acr IS NOT NULL) AND
                                                                                     (login_was_used IS NOT NULL) AND
                                                                                     (context IS NOT NULL) AND
                                                                                     (amr IS NOT NULL) AND
                                                                                     (consent_challenge_id IS NOT NULL) AND
                                                                                     (consent_verifier IS NOT NULL) AND
                                                                                     (consent_skip IS NOT NULL) AND
                                                                                     (consent_csrf IS NOT NULL) AND
                                                                                     (granted_scope IS NOT NULL) AND
                                                                                     (consent_remember IS NOT NULL) AND
                                                                                     (consent_remember_for IS NOT NULL) AND
                                                                                     (consent_error IS NOT NULL) AND
                                                                                     (session_access_token IS NOT NULL) AND
                                                                                     (session_id_token IS NOT NULL) AND
                                                                                     (consent_was_used IS NOT NULL))))),
  ALTER COLUMN requested_scope SET NOT NULL,
  ALTER COLUMN login_csrf SET NOT NULL,
  ALTER COLUMN subject SET NOT NULL,
  ALTER COLUMN request_url SET NOT NULL,
  ALTER COLUMN login_skip SET NOT NULL,
  ALTER COLUMN client_id SET NOT NULL,
  ALTER COLUMN oidc_context SET NOT NULL,
  ALTER COLUMN context SET NOT NULL,
  ALTER COLUMN state SET NOT NULL,
  ALTER COLUMN login_verifier SET NOT NULL,
  ALTER COLUMN login_remember SET NOT NULL,
  ALTER COLUMN login_remember_for SET NOT NULL,
  ALTER COLUMN acr SET NOT NULL,
  ALTER COLUMN login_was_used SET NOT NULL,
  ALTER COLUMN consent_skip SET NOT NULL,
  ALTER COLUMN consent_remember SET NOT NULL,
  ALTER COLUMN session_access_token SET NOT NULL,
  ALTER COLUMN session_id_token SET NOT NULL,
  ALTER COLUMN consent_was_used SET NOT NULL;
