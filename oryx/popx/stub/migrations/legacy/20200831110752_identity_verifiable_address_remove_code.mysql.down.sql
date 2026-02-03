ALTER TABLE `identity_verifiable_addresses` ADD COLUMN `code` VARCHAR (32);
ALTER TABLE `identity_verifiable_addresses` ADD COLUMN `expires_at` DATETIME;
UPDATE identity_verifiable_addresses SET code = LEFT(SHA2(RANDOM_BYTES(32), 256), 32) WHERE code IS NULL;
UPDATE identity_verifiable_addresses SET expires_at = CURRENT_TIMESTAMP WHERE expires_at IS NULL;
ALTER TABLE `identity_verifiable_addresses` MODIFY `code` VARCHAR (32) NOT NULL;
ALTER TABLE `identity_verifiable_addresses` MODIFY `expires_at` DATETIME;
CREATE UNIQUE INDEX `identity_verifiable_addresses_code_uq_idx` ON `identity_verifiable_addresses` (`code`);
CREATE INDEX `identity_verifiable_addresses_code_idx` ON `identity_verifiable_addresses` (`code`);
