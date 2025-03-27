ALTER TABLE hydra_client RENAME pk TO pk_deprecated;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
ALTER TABLE hydra_client ADD COLUMN pk UUID DEFAULT uuid_generate_v4();
ALTER TABLE hydra_client ALTER pk DROP DEFAULT;
ALTER TABLE hydra_client DROP CONSTRAINT hydra_client_pkey;
ALTER TABLE hydra_client ADD PRIMARY KEY (pk);
