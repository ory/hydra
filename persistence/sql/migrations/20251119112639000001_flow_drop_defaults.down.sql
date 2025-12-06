UPDATE hydra_oauth2_flow
SET login_extend_session_lifespan = COALESCE(login_extend_session_lifespan, FALSE),
    forced_subject_identifier     = COALESCE(forced_subject_identifier, '')
WHERE login_extend_session_lifespan IS NULL
   OR forced_subject_identifier IS NULL;
