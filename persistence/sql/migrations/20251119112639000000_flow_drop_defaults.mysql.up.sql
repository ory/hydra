ALTER TABLE hydra_oauth2_flow
  -- We need to drop these NOT NULL constraints, because the fields are actually not used anymore in the code, and therefore cannot be set by it.
  -- Mysql has issues with two changes to the same column in one statement, therefore we need to use MODIFY COLUMN here to both drop the default and constraint.
  MODIFY COLUMN forced_subject_identifier varchar(255) NULL,
  MODIFY COLUMN login_extend_session_lifespan tinyint(1) NULL,

  ALTER COLUMN requested_at_audience DROP DEFAULT,
  ALTER COLUMN amr DROP DEFAULT,
  ALTER COLUMN granted_at_audience DROP DEFAULT,
  ALTER COLUMN oidc_context DROP DEFAULT,
  ALTER COLUMN context DROP DEFAULT,
  ALTER COLUMN acr DROP DEFAULT,
  ALTER COLUMN consent_skip DROP DEFAULT,
  ALTER COLUMN consent_remember DROP DEFAULT,
  ALTER COLUMN login_remember DROP DEFAULT,
  ALTER COLUMN consent_was_used DROP DEFAULT,
  ALTER COLUMN login_was_used DROP DEFAULT,
  ALTER COLUMN session_id_token DROP DEFAULT,
  ALTER COLUMN session_access_token DROP DEFAULT;
