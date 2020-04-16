-- +migrate Up
CREATE UNIQUE INDEX hydra_oauth2_access_request_id_idx ON hydra_oauth2_access (request_id);
CREATE UNIQUE INDEX hydra_oauth2_refresh_request_id_idx ON hydra_oauth2_refresh (request_id);

-- +migrate Down
DROP INDEX hydra_oauth2_access_request_id_idx ON hydra_oauth2_access;
DROP INDEX hydra_oauth2_refresh_request_id_idx ON hydra_oauth2_refresh;
