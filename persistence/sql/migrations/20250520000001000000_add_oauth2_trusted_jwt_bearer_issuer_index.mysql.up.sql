-- We must avoid having empty fields at the beginning of an index, to avoid queries degenerating into full table scans.
-- Apart from this consideration, the order of the fields in the index is not really important when the query uses them all.
-- To avoid performance pitfalls/surprises, we place Ory-controlled fields before user-controlled fields.
CREATE UNIQUE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_nid_uq_idx ON hydra_oauth2_trusted_jwt_bearer_issuer (nid ASC, key_id ASC, issuer ASC, subject ASC);

