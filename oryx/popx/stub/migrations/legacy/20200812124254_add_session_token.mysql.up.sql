DELETE FROM sessions;
ALTER TABLE `sessions` ADD COLUMN `token` VARCHAR (32);
ALTER TABLE `sessions` MODIFY `token` VARCHAR (32);
CREATE UNIQUE INDEX `sessions_token_uq_idx` ON `sessions` (`token`);
CREATE INDEX `sessions_token_idx` ON `sessions` (`token`);