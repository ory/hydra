-- Supports the OR branch in the cascading consent revocation query for
-- RevokeSubjectConsentSession (subject = ? AND nid = ?) and
-- RevokeSubjectClientConsentSession (subject = ? AND client_id = ? AND nid = ?).
-- We place the Ory-controlled nid before the user-controlled subject so the
-- leading column is never empty; client_id trails so the same index also
-- serves the client-scoped revoke path as a perfect prefix match.
-- CockroachDB accepts CONCURRENTLY as a no-op (all CRDB indexes are built via
-- online schema change regardless); on Postgres it is required to avoid an
-- ACCESS EXCLUSIVE lock on large token tables.
CREATE INDEX CONCURRENTLY IF NOT EXISTS hydra_oauth2_access_nid_subject_idx ON hydra_oauth2_access (nid ASC, subject ASC, client_id ASC);
