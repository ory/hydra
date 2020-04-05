DELETE FROM hydra_oauth2_logout_request WHERE client_id IS NULL;
ALTER TABLE hydra_oauth2_logout_request MODIFY client_id varchar(255) NOT NULL;
