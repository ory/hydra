CREATE UNIQUE INDEX hydra_oauth2_logout_request_veri_idx ON hydra_oauth2_logout_request (verifier);
ALTER TABLE hydra_oauth2_logout_request ADD CONSTRAINT hydra_oauth2_logout_request_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;

