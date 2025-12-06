ALTER TABLE hydra_oauth2_flow
  ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP
    GENERATED ALWAYS AS (CASE WHEN consent_remember_for > 0 THEN requested_at + (consent_remember_for * INTERVAL '1 second') END) STORED; -- postgres supports virtual columns only with version 18+, so we have to use a stored column instead
