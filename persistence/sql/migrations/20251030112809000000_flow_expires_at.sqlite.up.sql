ALTER TABLE hydra_oauth2_flow
  ADD COLUMN expires_at TIMESTAMP
    GENERATED ALWAYS AS (if(consent_remember_for > 0, datetime(requested_at, '+' || consent_remember_for || ' seconds'), NULL)) VIRTUAL;
