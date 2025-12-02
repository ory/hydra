ALTER TABLE hydra_oauth2_flow
  MODIFY COLUMN login_extend_session_lifespan tinyint(1) NOT NULL DEFAULT FALSE,
  MODIFY COLUMN forced_subject_identifier varchar(255) NOT NULL DEFAULT '';
