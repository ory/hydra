ALTER TABLE hydra_client RENAME pk TO pk_deprecated;
-- UUID generation based on https://stackoverflow.com/a/21327318/12723442
ALTER TABLE hydra_client ADD COLUMN pk UUID DEFAULT uuid_in(
  overlay(
    overlay(
      substr(to_hex(random()::numeric || random()::numeric), 1, 32)
      placing '4'
      from 13
    )
    placing to_hex(floor(random() * 4 + 8)::int)::text
    from 17
  )::cstring
);
ALTER TABLE hydra_client ALTER pk DROP DEFAULT;
ALTER TABLE hydra_client DROP CONSTRAINT hydra_client_pkey;
ALTER TABLE hydra_client ADD PRIMARY KEY (pk);
