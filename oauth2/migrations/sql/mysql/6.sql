-- +migrate Up
CREATE INDEX hydra_oauth2_access_requested_at_idx ON hydra_oauth2_access (requested_at);

-- +migrate Down
DROP INDEX hydra_oauth2_access_requested_at_idx ON hydra_oauth2_access;
