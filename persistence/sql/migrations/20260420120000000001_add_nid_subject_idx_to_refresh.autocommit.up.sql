-- Companion to the access-table index; same rationale.
CREATE INDEX CONCURRENTLY IF NOT EXISTS hydra_oauth2_refresh_nid_subject_idx ON hydra_oauth2_refresh (nid ASC, subject ASC, client_id ASC);
