-- DROP INDEX hydra_oauth2_access_client_id_subject_idx ON hydra_oauth2_access;
DROP INDEX hydra_oauth2_access_expires_at_v2_idx ON hydra_oauth2_access; -- janitor still uses requested_at index

DROP INDEX hydra_oauth2_refresh_client_id_subject_idx ON hydra_oauth2_refresh;
DROP INDEX hydra_oauth2_refresh_expires_at_v2_idx ON hydra_oauth2_refresh; -- janitor still uses requested_at index

DROP INDEX hydra_oauth2_pkce_request_id_idx ON hydra_oauth2_pkce;
DROP INDEX hydra_oauth2_pkce_expires_at_v2_idx ON hydra_oauth2_pkce; -- janitor still uses requested_at index

DROP INDEX hydra_oauth2_oidc_request_id_idx ON hydra_oauth2_oidc;
DROP INDEX hydra_oauth2_oidc_expires_at_idx ON hydra_oauth2_oidc; -- janitor still uses requested_at index

DROP INDEX hydra_oauth2_code_request_id_idx ON hydra_oauth2_code;
DROP INDEX hydra_oauth2_code_expires_at_v2_idx ON hydra_oauth2_code; -- janitor still uses requested_at index
