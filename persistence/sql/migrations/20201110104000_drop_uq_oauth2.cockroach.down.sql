DROP INDEX hydra_oauth2_access_request_id_idx;
DROP INDEX hydra_oauth2_refresh_request_id_idx;

CREATE UNIQUE INDEX hydra_oauth2_access_request_id_idx ON hydra_oauth2_access (request_id);
CREATE UNIQUE INDEX hydra_oauth2_refresh_request_id_idx ON hydra_oauth2_refresh (request_id);

DROP INDEX hydra_oauth2_code_request_id_idx;
DROP INDEX hydra_oauth2_oidc_request_id_idx;
DROP INDEX hydra_oauth2_pkce_request_id_idx;
