ALTER TABLE hydra_oauth2_flow
  ADD COLUMN expires_at TIMESTAMP
    AS (IF(consent_remember_for > 0, DATE_ADD(requested_at, INTERVAL consent_remember_for SECOND), NULL)) VIRTUAL;
