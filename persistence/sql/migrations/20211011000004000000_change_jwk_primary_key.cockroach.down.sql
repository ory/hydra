-- No-op: the phantom index dropped by the up migration was erroneously
-- created by CockroachDB's schema changer and should not be recreated.

SELECT 1;
