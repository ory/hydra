CREATE TABLE IF NOT EXISTS hydra_oauth2_trusted_jwt_bearer_issuer
(
  id         UUID                    NOT NULL,
  issuer     VARCHAR(255)            NOT NULL,
  subject    VARCHAR(255)            NOT NULL,
  scope      TEXT                    NOT NULL,
  key_set    varchar(255)            NOT NULL,
  key_id     varchar(255)            NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  expires_at TIMESTAMP DEFAULT NOW() NOT NULL,
  UNIQUE (issuer, subject, key_id),
  FOREIGN KEY (key_set, key_id) REFERENCES hydra_jwk (sid, kid) ON DELETE CASCADE,
  CONSTRAINT "primary" PRIMARY KEY (id ASC)
);

CREATE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_expires_at_idx ON hydra_oauth2_trusted_jwt_bearer_issuer (expires_at);
