-- +migrate Up
CREATE INDEX hydra_oauth2_logout_request_client_id_idx ON hydra_oauth2_logout_request (client_id);

-- +migrate Down
ALTER TABLE hydra_oauth2_logout_request DROP FOREIGN KEY hydra_oauth2_logout_request_client_id_fk;
DROP INDEX hydra_oauth2_logout_request_client_id_idx ON hydra_oauth2_logout_request;
ALTER TABLE hydra_oauth2_logout_request ADD CONSTRAINT hydra_oauth2_logout_request_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
