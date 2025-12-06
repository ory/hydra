ALTER TABLE hydra_oauth2_flow
  ALTER COLUMN forced_subject_identifier SET DEFAULT '',
  ALTER COLUMN login_extend_session_lifespan SET DEFAULT FALSE,

  MODIFY COLUMN requested_at_audience JSON DEFAULT (_utf8mb4'[]'),
  MODIFY COLUMN amr JSON DEFAULT (_utf8mb4'[]'),
  MODIFY COLUMN granted_at_audience JSON DEFAULT (_utf8mb4'[]'),
  MODIFY COLUMN oidc_context JSON NOT NULL DEFAULT (_utf8mb4'{}'),
  MODIFY COLUMN context JSON NOT NULL DEFAULT (_utf8mb4'{}'),
  MODIFY COLUMN acr TEXT NOT NULL DEFAULT (_utf8mb4''),
  ALTER COLUMN consent_skip SET DEFAULT FALSE,
  ALTER COLUMN consent_remember SET DEFAULT FALSE,
  ALTER COLUMN login_remember SET DEFAULT FALSE,
  ALTER COLUMN consent_was_used SET DEFAULT FALSE,
  ALTER COLUMN login_was_used SET DEFAULT FALSE,
  MODIFY COLUMN session_id_token JSON NOT NULL DEFAULT (_utf8mb4'{}'),
  MODIFY COLUMN session_access_token JSON NOT NULL DEFAULT (_utf8mb4'{}');
