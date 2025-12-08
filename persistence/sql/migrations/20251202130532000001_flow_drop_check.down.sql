DELETE
FROM hydra_oauth2_flow
WHERE requested_scope IS NULL
   OR login_csrf IS NULL
   OR subject IS NULL
   OR request_url IS NULL
   OR login_skip IS NULL
   OR client_id IS NULL
   OR oidc_context IS NULL
   OR context IS NULL
   OR state IS NULL
   OR login_verifier IS NULL
   OR login_remember IS NULL
   OR login_remember_for IS NULL
   OR acr IS NULL
   OR login_was_used IS NULL
   OR consent_skip IS NULL
   OR consent_remember IS NULL
   OR session_access_token IS NULL
   OR session_id_token IS NULL
   OR consent_was_used IS NULL;
