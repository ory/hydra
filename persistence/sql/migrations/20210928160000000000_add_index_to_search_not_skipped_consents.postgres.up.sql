CREATE INDEX IF NOT EXISTS hydra_oauth2_consent_request_client_id_subject_not_skipped
    ON hydra_oauth2_consent_request (client_id, subject, skip)
    WHERE NOT skip;
