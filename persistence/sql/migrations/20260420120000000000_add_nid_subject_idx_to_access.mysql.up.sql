-- MySQL variant: columns ordered (subject, nid, client_id) so this index does
-- NOT back the FOREIGN KEY (nid) REFERENCES networks(id) constraint. MySQL
-- auto-binds any FK-leading-column index as the FK's sole backing, which then
-- blocks drops. Equality predicates on (subject, nid) and
-- (subject, nid, client_id) are served optimally either way.
CREATE INDEX hydra_oauth2_access_subject_nid_idx ON hydra_oauth2_access (subject ASC, nid ASC, client_id ASC);
