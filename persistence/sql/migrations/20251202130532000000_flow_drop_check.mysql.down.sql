-- this is not ideal, but required because of MySQL limitations regarding changing columns that are used in foreign key constraints
SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE hydra_oauth2_flow
  ADD CONSTRAINT hydra_oauth2_flow_chk CHECK (((state = 128) OR (state = 129) OR (state = 1) OR ((state = 2) AND
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

  MODIFY COLUMN requested_scope json NOT NULL,
  MODIFY COLUMN login_csrf VARCHAR (40) NOT NULL,
  MODIFY COLUMN subject VARCHAR (255) NOT NULL,
  MODIFY COLUMN request_url TEXT NOT NULL,
  MODIFY COLUMN login_skip tinyint(1) NOT NULL,
  MODIFY COLUMN client_id varchar(255) NOT NULL,
  MODIFY COLUMN oidc_context json NOT NULL,
  MODIFY COLUMN context json NOT NULL,
  MODIFY COLUMN state SMALLINT NOT NULL,
  MODIFY COLUMN acr TEXT NOT NULL,
  MODIFY COLUMN consent_skip tinyint(1) NOT NULL,
  MODIFY COLUMN consent_remember tinyint(1) NOT NULL,
  MODIFY COLUMN login_remember tinyint(1) NOT NULL,
  MODIFY COLUMN consent_was_used tinyint(1) NOT NULL,
  MODIFY COLUMN login_was_used tinyint(1) NOT NULL,
  MODIFY COLUMN session_id_token json NOT NULL,
  MODIFY COLUMN session_access_token json NOT NULL;

SET FOREIGN_KEY_CHECKS = 1;
