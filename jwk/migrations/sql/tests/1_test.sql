-- +migrate Up
INSERT INTO hydra_jwk (sid, kid, version, keydata) VALUES ('1-sid', '1-kid', 0, 'some-key');

-- +migrate Down
