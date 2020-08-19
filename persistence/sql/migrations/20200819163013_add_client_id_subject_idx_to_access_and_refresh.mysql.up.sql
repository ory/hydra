CREATE INDEX hydra_oauth2_access_client_id_subject_idx ON hydra_oauth2_access (client_id, subject);
CREATE INDEX hydra_oauth2_refresh_client_id_subject_idx ON hydra_oauth2_refresh (client_id, subject);
