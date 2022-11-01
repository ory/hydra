SET @dbName = DATABASE();

DROP FUNCTION IF EXISTS isFieldExisting;
DROP PROCEDURE IF EXISTS addFieldIfNotExists;

SET GLOBAL log_bin_trust_function_creators = 1;

CREATE FUNCTION isFieldExisting (
  dbName VARCHAR(255),
  tableName VARCHAR(255),
  columnName VARCHAR(255)
)
  RETURNS INT
  RETURN (
    SELECT COUNT(COLUMN_NAME)
    FROM INFORMATION_SCHEMA.columns
    WHERE TABLE_SCHEMA = dbName
      AND TABLE_NAME = tableName
      AND COLUMN_NAME = columnName
    LIMIT 1
);

CREATE PROCEDURE addFieldIfNotExists(
  IN  dbName varchar(255),
  IN  tableName varchar(255),
  IN  columnName varchar(255),
  IN  definition varchar(255)
)
BEGIN

  SET @isFieldThere = isFieldExisting(dbName, tableName, columnName);
  IF (@isFieldThere = 0) THEN

    SET @ddl = CONCAT('ALTER TABLE ', tableName);
    SET @ddl = CONCAT(@ddl, ' ', 'ADD COLUMN') ;
    SET @ddl = CONCAT(@ddl, ' ', columnName);
    SET @ddl = CONCAT(@ddl, ' ', definition);

    PREPARE stmt FROM @ddl;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

  END IF;
END;

SET GLOBAL log_bin_trust_function_creators = 0;

CALL addFieldIfNotExists(@dbName, 'hydra_oauth2_consent_request', 'amr', 'TEXT NULL');
CALL addFieldIfNotExists(@dbName, 'hydra_oauth2_authentication_request_handled', 'amr', 'TEXT NULL');

DROP FUNCTION IF EXISTS isFieldExisting;
DROP PROCEDURE IF EXISTS addFieldIfNotExists;

UPDATE hydra_oauth2_consent_request SET amr='' WHERE amr IS NULL;
ALTER TABLE hydra_oauth2_consent_request MODIFY amr TEXT NOT NULL;

UPDATE hydra_oauth2_authentication_request_handled SET amr='' WHERE amr IS NULL;
ALTER TABLE hydra_oauth2_authentication_request_handled MODIFY amr TEXT NOT NULL;
