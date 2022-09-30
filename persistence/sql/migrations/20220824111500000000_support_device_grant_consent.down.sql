ALTER TABLE
    hydra_oauth2_device_grant_request DROP FOREIGN KEY IF EXISTS hydra_oauth2_device_grant_request_client_id_fk;

ALTER TABLE
    hydra_oauth2_device_grant_request DROP FOREIGN KEY IF EXISTS hydra_oauth2_device_grant_request_nid_fk_idx;

DROP TABLE IF EXISTS hydra_oauth2_device_grant_request;