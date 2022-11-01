ALTER TABLE hydra_oauth2_consent_request DROP COLUMN amr;
ALTER TABLE hydra_oauth2_authentication_request_handled DROP COLUMN amr;

DROP FUNCTION IF EXISTS isFieldExisting;
DROP PROCEDURE IF EXISTS addFieldIfNotExists;

