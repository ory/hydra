DELETE FROM identity_recovery_tokens WHERE selfservice_recovery_flow_id IS NULL;
ALTER TABLE `identity_recovery_tokens` MODIFY `selfservice_recovery_flow_id` char(36) NOT NULL;
ALTER TABLE `identity_recovery_tokens` DROP COLUMN `expires_at`;
ALTER TABLE `identity_recovery_tokens` DROP COLUMN `issued_at`;