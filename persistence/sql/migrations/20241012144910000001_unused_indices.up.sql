DROP INDEX IF EXISTS hydra_oauth2_access_client_id_subject_idx;
DROP INDEX IF EXISTS hydra_oauth2_access_expires_at_v2_idx; -- janitor still uses requested_at index

DROP INDEX IF EXISTS hydra_oauth2_refresh_client_id_subject_idx;
DROP INDEX IF EXISTS hydra_oauth2_refresh_expires_at_v2_idx; -- janitor still uses requested_at index

DROP INDEX IF EXISTS hydra_oauth2_pkce_request_id_idx;
DROP INDEX IF EXISTS hydra_oauth2_pkce_expires_at_v2_idx; -- janitor still uses requested_at index

DROP INDEX IF EXISTS hydra_oauth2_oidc_request_id_idx;
DROP INDEX IF EXISTS hydra_oauth2_oidc_expires_at_idx; -- janitor still uses requested_at index

DROP INDEX IF EXISTS hydra_oauth2_code_request_id_idx;
DROP INDEX IF EXISTS hydra_oauth2_code_expires_at_v2_idx; -- janitor still uses requested_at index
