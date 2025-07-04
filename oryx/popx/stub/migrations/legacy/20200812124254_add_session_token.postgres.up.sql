DELETE FROM sessions;
ALTER TABLE "sessions" ADD COLUMN "token" VARCHAR (32);
ALTER TABLE "sessions" ALTER COLUMN "token" TYPE VARCHAR (32), ALTER COLUMN "token" DROP NOT NULL;
CREATE UNIQUE INDEX "sessions_token_uq_idx" ON "sessions" (token);
CREATE INDEX "sessions_token_idx" ON "sessions" (token);