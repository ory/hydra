-- +migrate Up
CREATE TABLE IF NOT EXISTS hydra_client (
  pk SERIAL PRIMARY KEY,
	id varchar(255) NOT NULL,
	client_name text NOT NULL,
	client_secret text NOT NULL,
	redirect_uris text NOT NULL,
	grant_types text NOT NULL,
	response_types text NOT NULL,
	scope text NOT NULL,
	owner text NOT NULL,
	policy_uri text NOT NULL,
	tos_uri text NOT NULL,
	client_uri text NOT NULL,
	logo_uri text NOT NULL,
	contacts text NOT NULL,
	client_secret_expires_at INTEGER NOT NULL DEFAULT 0,
	sector_identifier_uri text NOT NULL,
	jwks text NOT NULL,
	jwks_uri text NOT NULL,
	request_uris text NOT NULL,
	token_endpoint_auth_method VARCHAR(25) NOT NULL DEFAULT '',
	request_object_signing_alg  VARCHAR(10) NOT NULL DEFAULT '',
	userinfo_signed_response_alg VARCHAR(10) NOT NULL DEFAULT '',
	subject_type VARCHAR(15) NOT NULL DEFAULT '',
	allowed_cors_origins text NOT NULL,
	audience text NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	UNIQUE (id)
);

-- +migrate Down
DROP TABLE hydra_client;
