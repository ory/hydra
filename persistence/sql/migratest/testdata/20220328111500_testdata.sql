INSERT INTO hydra_jwk (pk_deprecated, pk, sid, kid, nid, version, keydata, created_at) VALUES (9, '98565339-57c7-4bc0-bc3d-53171d60e832', 'sid-0009', 'kid-0009', (SELECT id FROM networks LIMIT 1), 2, 'key-0002', CURRENT_TIMESTAMP);

INSERT INTO hydra_oauth2_trusted_jwt_bearer_issuer (id, nid, issuer, subject, allow_any_subject, scope, key_set, key_id)
VALUES ('30e51720-4a88-48ca-8243-de7d8f461675', (SELECT id FROM networks LIMIT 1), 'some-issuer', 'some-subject', false, 'some-scope', 'sid-0009', 'kid-0009');

INSERT INTO hydra_oauth2_trusted_jwt_bearer_issuer (id, nid, issuer, subject, allow_any_subject, scope, key_set, key_id)
VALUES ('30e51720-4a88-48ca-8243-de7d8f461676', (SELECT id FROM networks LIMIT 1), 'some-issuer', '', true, 'some-scope', 'sid-0009', 'kid-0009');
