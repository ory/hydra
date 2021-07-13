-- +migrate Up
DELETE FROM hydra_jwk WHERE sid='hydra.openid.id-token';

-- +migrate Down
