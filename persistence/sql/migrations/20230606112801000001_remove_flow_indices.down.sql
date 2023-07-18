CREATE UNIQUE INDEX hydra_oauth2_flow_login_verifier_idx ON hydra_oauth2_flow (login_verifier);
CREATE UNIQUE INDEX hydra_oauth2_flow_consent_verifier_idx ON hydra_oauth2_flow (consent_verifier);

CREATE INDEX hydra_oauth2_flow_multi_query_idx
  ON hydra_oauth2_flow
    (
     consent_error ASC, state ASC, subject ASC,
     client_id ASC, consent_skip ASC, consent_remember
     ASC, nid ASC
      );

DROP INDEX hydra_oauth2_flow_previous_consents_idx;
