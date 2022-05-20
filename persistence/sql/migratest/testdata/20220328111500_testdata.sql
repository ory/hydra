INSERT INTO hydra_jwk (pk, sid, kid, version, keydata, created_at) VALUES (9, 'sid-0009', 'kid-0009', 2, 'key-0002', now());

INSERT INTO hydra_oauth2_trusted_jwt_bearer_issuer (id, issuer, subject, allow_any_subject, scope, key_set, key_id)
VALUES ('30e51720-4a88-48ca-8243-de7d8f461675', 'some-issuer', 'some-subject', false, 'some-scope', 'sid-0009', 'kid-0009');

INSERT INTO hydra_oauth2_trusted_jwt_bearer_issuer (id, issuer, subject, allow_any_subject, scope, key_set, key_id)
VALUES ('30e51720-4a88-48ca-8243-de7d8f461676', 'some-issuer', '', true, 'some-scope', 'sid-0009', 'kid-0009');
