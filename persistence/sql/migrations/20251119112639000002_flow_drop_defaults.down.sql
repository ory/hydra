ALTER TABLE hydra_oauth2_flow
  ALTER COLUMN login_extend_session_lifespan SET DEFAULT FALSE,
  ALTER COLUMN forced_subject_identifier SET DEFAULT '',

  ALTER COLUMN requested_at_audience SET DEFAULT '[]'::jsonb,
  ALTER COLUMN oidc_context SET DEFAULT '{}'::jsonb,
  ALTER COLUMN context SET DEFAULT '{}'::jsonb,
  ALTER COLUMN amr SET DEFAULT '[]'::jsonb,
  ALTER COLUMN acr SET DEFAULT '',
  ALTER COLUMN consent_skip SET DEFAULT FALSE,
  ALTER COLUMN granted_at_audience SET DEFAULT '[]'::jsonb,
  ALTER COLUMN consent_remember SET DEFAULT FALSE,
  ALTER COLUMN login_remember SET DEFAULT FALSE,
  ALTER COLUMN consent_was_used SET DEFAULT FALSE,
  ALTER COLUMN login_was_used SET DEFAULT FALSE,
  ALTER COLUMN session_id_token SET DEFAULT '{}'::jsonb,
  ALTER COLUMN session_access_token SET DEFAULT '{}'::jsonb;
