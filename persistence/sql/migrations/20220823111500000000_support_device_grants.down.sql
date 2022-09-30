ALTER TABLE
    hydra_oauth2_device_code DROP FOREIGN KEY IF EXISTS hydra_oauth2_device_code_challenge_id_fk;

ALTER TABLE
    hydra_oauth2_device_code DROP FOREIGN KEY IF EXISTS hydra_oauth2_device_code_client_id_fk;

ALTER TABLE
    hydra_oauth2_device_code DROP FOREIGN KEY IF EXISTS hydra_oauth2_device_code_nid_fk_idx;

DROP TABLE IF EXISTS hydra_oauth2_device_code;

ALTER TABLE
    hydra_oauth2_user_code DROP FOREIGN KEY IF EXISTS hydra_oauth2_user_code_challenge_id_fk;

ALTER TABLE
    hydra_oauth2_user_code DROP FOREIGN KEY IF EXISTS hydra_oauth2_user_code_client_id_fk;

ALTER TABLE
    hydra_oauth2_user_code DROP FOREIGN KEY IF EXISTS hydra_oauth2_user_code_nid_fk_idx;

DROP TABLE IF EXISTS hydra_oauth2_user_code;