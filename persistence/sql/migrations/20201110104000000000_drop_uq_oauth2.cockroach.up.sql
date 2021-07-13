--- This is hardcoded and if it fails we need to manually figure out the constraint name:
--- SELECT constraint_name FROM information_schema.table_constraints WHERE table_name='hydra_oauth2_access' AND constraint_type='UNIQUE';
DROP INDEX hydra_oauth2_access_request_id_key CASCADE;
DROP INDEX hydra_oauth2_refresh_request_id_key CASCADE;

CREATE INDEX hydra_oauth2_access_request_id_idx ON hydra_oauth2_access (request_id);
CREATE INDEX hydra_oauth2_refresh_request_id_idx ON hydra_oauth2_refresh (request_id);

CREATE INDEX hydra_oauth2_code_request_id_idx ON hydra_oauth2_code (request_id);
CREATE INDEX hydra_oauth2_oidc_request_id_idx ON hydra_oauth2_oidc (request_id);
CREATE INDEX hydra_oauth2_pkce_request_id_idx ON hydra_oauth2_pkce (request_id);
