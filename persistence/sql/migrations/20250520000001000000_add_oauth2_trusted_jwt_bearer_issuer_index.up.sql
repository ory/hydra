-- `key_id` is unique-ish per row so we place it first in the index to make queries including it very fast.
CREATE UNIQUE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_key_id ON hydra_oauth2_trusted_jwt_bearer_issuer (key_id ASC, issuer ASC, subject ASC, nid ASC);
