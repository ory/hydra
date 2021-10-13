CREATE TABLE "_hydra_jwk_tmp"
(
    sid             VARCHAR(255) NOT NULL,
    kid             VARCHAR(255) NOT NULL,
    version         INTEGER      DEFAULT 0 NOT NULL,
    keydata         TEXT         NOT NULL,
    created_at      TIMESTAMP    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    pk              TEXT         PRIMARY KEY,
    pk_deprecated   INTEGER      NULL DEFAULT NULL
);

-- UUID generation based on https://stackoverflow.com/a/61000724/12723442
INSERT INTO "_hydra_jwk_tmp" (
    sid,
    kid,
    version,
    keydata,
    created_at,
    pk,
    pk_deprecated
) SELECT
    sid,
    kid,
    version,
    keydata,
    created_at,
    lower(
      hex(randomblob(4)) ||
      '-' || hex(randomblob(2)) ||
      '-' || '4' || substr(hex(randomblob(2)), 2) ||
      '-' || substr('AB89', 1 + (abs(random()) % 4) , 1) || substr(hex(randomblob(2)), 2) ||
      '-' || hex(randomblob(6))
    ),
    pk
FROM "hydra_jwk";

DROP INDEX hydra_jwk_sid_kid_key;
DROP TABLE "hydra_jwk";
ALTER TABLE "_hydra_jwk_tmp" RENAME TO "hydra_jwk";

CREATE UNIQUE INDEX hydra_jwk_sid_kid_key ON hydra_jwk (sid, kid);
