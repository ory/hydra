DROP INDEX `identity_verifiable_addresses_code_uq_idx` ON `identity_verifiable_addresses`;
DROP INDEX `identity_verifiable_addresses_code_idx` ON `identity_verifiable_addresses`;
ALTER TABLE `identity_verifiable_addresses` DROP COLUMN `code`;
ALTER TABLE `identity_verifiable_addresses` DROP COLUMN `expires_at`;