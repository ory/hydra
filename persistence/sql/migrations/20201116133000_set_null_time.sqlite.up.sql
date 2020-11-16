CREATE TABLE "_hydra_oauth2_authentication_session_tmp"
(
    id               VARCHAR(40)  NOT NULL PRIMARY KEY,
    authenticated_at TIMESTAMP    NULL,
    subject          VARCHAR(255) NOT NULL,
    remember         INTEGER      NOT NULL DEFAULT false
);

INSERT INTO "_hydra_oauth2_authentication_session_tmp" (
    id, authenticated_at, subject, remember
) SELECT id, authenticated_at, subject, remember FROM "hydra_oauth2_authentication_session";

DROP INDEX hydra_oauth2_authentication_session_subject_idx;
DROP TABLE "hydra_oauth2_authentication_session";
ALTER TABLE "_hydra_oauth2_authentication_session_tmp" RENAME TO "hydra_oauth2_authentication_session";

CREATE INDEX hydra_oauth2_authentication_session_subject_idx ON hydra_oauth2_authentication_session (subject);
