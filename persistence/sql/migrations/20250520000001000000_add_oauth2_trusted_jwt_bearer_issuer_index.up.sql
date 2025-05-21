-- `key_id` is unique-ish per row so we place it first in the index to make queries including it very fast.
-- Other fields have very few distinct values in the table so having them first in the index makes queries do a full table scan.
CREATE UNIQUE INDEX IF NOT EXISTS hydra_oauth2_trusted_jwt_bearer_issuer_key_id ON hydra_oauth2_trusted_jwt_bearer_issuer (key_id ASC, issuer ASC, subject ASC, nid ASC);
