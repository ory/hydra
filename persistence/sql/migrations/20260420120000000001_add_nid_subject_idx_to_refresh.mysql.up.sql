-- MySQL variant: see companion access migration for the column-order rationale.
CREATE INDEX hydra_oauth2_refresh_subject_nid_idx ON hydra_oauth2_refresh (subject ASC, nid ASC, client_id ASC);
