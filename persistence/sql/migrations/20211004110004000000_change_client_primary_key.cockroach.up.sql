-- CockroachDB's declarative schema changer (observed on v26.1.1 and v26.1.2)
-- creates a phantom secondary UNIQUE index on the old primary key columns when
-- using `ALTER TABLE ... DROP CONSTRAINT "primary", ADD CONSTRAINT ... PRIMARY KEY`,
-- contradicting the documented behavior that this pattern replaces the primary
-- key without creating a secondary index. This step drops the phantom index
-- left behind by the step-3 migration above.
--
-- The client.Client struct has never had a pk_deprecated field, so this
-- phantom index has been inert. We drop it as defensive cleanup: it is a
-- schema artifact that should not exist, and any future regression that
-- introduces a constant pk_deprecated write would reactivate the same bug
-- that caused production 409 Conflict errors on hydra_jwk.
--
-- IF EXISTS is used because customers whose migrations ran on older
-- CockroachDB versions (where the legacy schema changer honored the
-- documented behavior) will not have this phantom index.

DROP INDEX IF EXISTS hydra_client_pk_deprecated_key CASCADE;
