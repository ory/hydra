UPDATE selfservice_errors SET seen_at = '1980-01-01 00:00:00' WHERE seen_at = NULL;
ALTER TABLE `selfservice_errors` MODIFY `seen_at` DATETIME;