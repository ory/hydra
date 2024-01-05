CREATE INDEX hydra_jwk_nid_sid_created_at_idx ON hydra_jwk (nid, sid, created_at);
CREATE INDEX hydra_jwk_nid_sid_kid_created_at_idx ON hydra_jwk (nid, sid, kid, created_at);
