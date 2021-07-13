-- +migrate Up
ALTER TABLE hydra_oauth2_access ADD subject varchar(255) NOT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_refresh ADD subject varchar(255) NOT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_code ADD subject varchar(255) NOT NULL DEFAULT '';
ALTER TABLE hydra_oauth2_oidc ADD subject varchar(255) NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE hydra_oauth2_access  DROP COLUMN subject;
ALTER TABLE hydra_oauth2_refresh DROP COLUMN subject;
ALTER TABLE hydra_oauth2_code DROP COLUMN subject;
ALTER TABLE hydra_oauth2_oidc DROP COLUMN subject;
