ALTER TABLE hydra_oauth2_flow
  ALTER COLUMN login_extend_session_lifespan SET NOT NULL,
  ALTER COLUMN forced_subject_identifier SET NOT NULL;
