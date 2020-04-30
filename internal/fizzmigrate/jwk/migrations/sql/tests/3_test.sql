-- +migrate Up
INSERT INTO hydra_jwk (sid, kid, version, keydata, created_at) VALUES ('3-sid', '3-kid', 0, 'some-key', NOW());

-- +migrate Down
