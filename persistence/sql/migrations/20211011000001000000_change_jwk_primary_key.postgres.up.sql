ALTER TABLE hydra_jwk RENAME pk TO pk_deprecated;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
ALTER TABLE hydra_jwk ADD COLUMN pk UUID DEFAULT uuid_generate_v4();
ALTER TABLE hydra_jwk ALTER pk DROP DEFAULT;
ALTER TABLE hydra_jwk DROP CONSTRAINT hydra_jwk_pkey;
ALTER TABLE hydra_jwk ADD PRIMARY KEY (pk);
