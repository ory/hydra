-- +migrate Up
CREATE INDEX hydra_oauth2_logout_request_client_id_idx ON hydra_oauth2_logout_request (client_id);

-- +migrate Down
DROP INDEX hydra_oauth2_logout_request_client_id_idx;
