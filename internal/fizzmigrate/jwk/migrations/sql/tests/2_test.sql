-- +migrate Up
INSERT INTO hydra_jwk (sid, kid, version, keydata, created_at) VALUES ('2-sid', '2-kid', 0, 'some-key', NOW());

-- +migrate Down
