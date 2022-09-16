
CREATE INDEX hydra_oauth2_flow_consent_error_state_subject_client_id_consent ON hydra_oauth2_flow USING BTREE (consent_error ASC, state ASC, subject ASC, client_id ASC, consent_skip ASC, consent_remember ASC, nid ASC);
