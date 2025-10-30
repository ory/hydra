ALTER TABLE hydra_oauth2_flow
  ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP
    AS (IF(consent_remember_for > 0, requested_at + INTERVAL '1 second' * consent_remember_for, NULL)) VIRTUAL;
