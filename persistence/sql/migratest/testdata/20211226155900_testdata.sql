INSERT INTO hydra_jwk (pk, sid, kid, version, keydata, created_at) VALUES ('e18d8447-3ec2-42d9-a3ad-e7cca8aa81f0', 'sid-0008', 'kid-0008', 2, 'key-0002', '2022-02-15 22:20:23');

INSERT INTO hydra_oauth2_trusted_jwt_bearer_issuer (id, issuer, subject, scope, key_set, key_id)
VALUES ('30e51720-4a88-48ca-8243-de7d8f461674', 'some-issuer', 'some-subject', 'some-scope', 'sid-0008', 'kid-0008');
