CREATE TABLE "_hydra_jwk_tmp"
(
    sid             VARCHAR(255) NOT NULL,
    kid             VARCHAR(255) NOT NULL,
    version         INTEGER      DEFAULT 0 NOT NULL,
    keydata         TEXT         NOT NULL,
    created_at      TIMESTAMP    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    pk              INTEGER      PRIMARY KEY
);

INSERT INTO "_hydra_jwk_tmp" (
    sid,
    kid,
    version,
    keydata,
    created_at,
    pk
) SELECT
    sid,
    kid,
    version,
    keydata,
    created_at,
    pk_deprecated
FROM "hydra_jwk";

DROP INDEX hydra_jwk_sid_kid_key;
DROP TABLE "hydra_jwk";
ALTER TABLE "_hydra_jwk_tmp" RENAME TO "hydra_jwk";

CREATE UNIQUE INDEX hydra_jwk_sid_kid_key ON hydra_jwk (sid, kid);
