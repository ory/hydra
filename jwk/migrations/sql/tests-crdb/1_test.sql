-- +migrate Up
INSERT INTO hydra_jwk (sid, kid, version, keydata, created_at) VALUES ('4-sid', '4-kid', 0, 'some-key', NOW());

-- +migrate Down
