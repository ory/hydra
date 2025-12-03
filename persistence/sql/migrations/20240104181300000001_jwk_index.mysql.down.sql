DROP INDEX hydra_jwk_nid_sid_created_at_idx ON hydra_jwk;
-- can't drop this because of a foreign key constraint
-- DROP INDEX hydra_jwk_nid_sid_kid_created_at_idx ON hydra_jwk;
