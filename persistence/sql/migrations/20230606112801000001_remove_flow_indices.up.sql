DROP INDEX hydra_oauth2_flow_login_verifier_idx;
DROP INDEX hydra_oauth2_flow_consent_verifier_idx;
DROP INDEX hydra_oauth2_flow_multi_query_idx;

CREATE INDEX IF NOT EXISTS hydra_oauth2_flow_previous_consents_idx
  ON hydra_oauth2_flow (subject, client_id, nid, consent_skip, consent_error, consent_remember);
