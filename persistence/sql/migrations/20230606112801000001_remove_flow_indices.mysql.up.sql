DROP INDEX hydra_oauth2_flow_login_verifier_idx ON hydra_oauth2_flow;
DROP INDEX hydra_oauth2_flow_consent_verifier_idx ON hydra_oauth2_flow;
DROP INDEX hydra_oauth2_flow_multi_query_idx ON hydra_oauth2_flow;

CREATE INDEX hydra_oauth2_flow_previous_consents_idx
  ON hydra_oauth2_flow (subject, client_id, nid, consent_skip, consent_error(2), consent_remember);
