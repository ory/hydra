-- migrations hash: f5af1bbf8624fd09cf094c1a5745a255e6ea8d56bc7177e0db9eb73d72f1b3dae0fdf3312458a84550c1f2374b0a1ea6fdd026e22267fbf1808b8e8611fb39c0

CREATE TABLE public.schema_migration (
	version VARCHAR(48) NOT NULL,
	version_self INT8 NOT NULL DEFAULT 0:::INT8,
	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
	CONSTRAINT schema_migration_pkey PRIMARY KEY (rowid ASC),
	UNIQUE INDEX schema_migration_version_idx (version ASC),
	INDEX schema_migration_version_self_idx (version_self ASC)
);
CREATE TABLE public.networks (
	id UUID NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	CONSTRAINT networks_pkey PRIMARY KEY (id ASC)
);
CREATE TABLE public.hydra_client (
	id VARCHAR(255) NOT NULL,
	client_name STRING NOT NULL,
	client_secret STRING NOT NULL,
	scope STRING NOT NULL,
	owner STRING NOT NULL,
	policy_uri STRING NOT NULL,
	tos_uri STRING NOT NULL,
	client_uri STRING NOT NULL,
	logo_uri STRING NOT NULL,
	client_secret_expires_at INT8 NOT NULL DEFAULT 0:::INT8,
	sector_identifier_uri STRING NOT NULL,
	jwks STRING NOT NULL,
	jwks_uri STRING NOT NULL,
	token_endpoint_auth_method VARCHAR(25) NOT NULL DEFAULT '':::STRING,
	request_object_signing_alg VARCHAR(10) NOT NULL DEFAULT '':::STRING,
	userinfo_signed_response_alg VARCHAR(10) NOT NULL DEFAULT '':::STRING,
	subject_type VARCHAR(15) NOT NULL DEFAULT '':::STRING,
	pk_deprecated INT8 NOT NULL DEFAULT unique_rowid(),
	created_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	frontchannel_logout_uri STRING NOT NULL DEFAULT '':::STRING,
	frontchannel_logout_session_required BOOL NOT NULL DEFAULT false,
	backchannel_logout_uri STRING NOT NULL DEFAULT '':::STRING,
	backchannel_logout_session_required BOOL NOT NULL DEFAULT false,
	metadata STRING NOT NULL DEFAULT '{}':::STRING,
	token_endpoint_auth_signing_alg VARCHAR(10) NOT NULL DEFAULT '':::STRING,
	authorization_code_grant_access_token_lifespan INT8 NULL,
	authorization_code_grant_id_token_lifespan INT8 NULL,
	authorization_code_grant_refresh_token_lifespan INT8 NULL,
	client_credentials_grant_access_token_lifespan INT8 NULL,
	implicit_grant_access_token_lifespan INT8 NULL,
	implicit_grant_id_token_lifespan INT8 NULL,
	jwt_bearer_grant_access_token_lifespan INT8 NULL,
	password_grant_access_token_lifespan INT8 NULL,
	password_grant_refresh_token_lifespan INT8 NULL,
	refresh_token_grant_id_token_lifespan INT8 NULL,
	refresh_token_grant_access_token_lifespan INT8 NULL,
	refresh_token_grant_refresh_token_lifespan INT8 NULL,
	pk UUID NULL,
	registration_access_token_signature VARCHAR(128) NOT NULL DEFAULT '':::STRING,
	nid UUID NOT NULL,
	redirect_uris JSONB NOT NULL,
	grant_types JSONB NOT NULL,
	response_types JSONB NOT NULL,
	audience JSONB NOT NULL,
	allowed_cors_origins JSONB NOT NULL,
	contacts JSONB NOT NULL,
	request_uris JSONB NOT NULL,
	post_logout_redirect_uris JSONB NOT NULL DEFAULT '[]':::JSONB,
	access_token_strategy VARCHAR(10) NOT NULL DEFAULT '':::STRING,
	skip_consent BOOL NOT NULL DEFAULT false,
	skip_logout_consent BOOL NULL,
	device_authorization_grant_id_token_lifespan INT8 NULL,
	device_authorization_grant_access_token_lifespan INT8 NULL,
	device_authorization_grant_refresh_token_lifespan INT8 NULL,
	CONSTRAINT hydra_client_pkey PRIMARY KEY (id ASC, nid ASC),
	UNIQUE INDEX hydra_client_id_key (id ASC, nid ASC),
	UNIQUE INDEX hydra_client_pk_key (pk ASC)
);
CREATE TABLE public.hydra_jwk (
	sid VARCHAR(255) NOT NULL,
	kid VARCHAR(255) NOT NULL,
	version INT8 NOT NULL DEFAULT 0:::INT8,
	keydata STRING NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	pk_deprecated INT8 NOT NULL DEFAULT unique_rowid(),
	pk UUID NOT NULL,
	nid UUID NOT NULL,
	CONSTRAINT hydra_jwk_pkey PRIMARY KEY (pk ASC),
	UNIQUE INDEX hydra_jwk_sid_kid_nid_key (sid ASC, kid ASC, nid ASC),
	INDEX hydra_jwk_nid_sid_created_at_idx (nid ASC, sid ASC, created_at ASC),
	INDEX hydra_jwk_nid_sid_kid_created_at_idx (nid ASC, sid ASC, kid ASC, created_at ASC)
);
CREATE TABLE public.hydra_oauth2_authentication_session (
	id VARCHAR(40) NOT NULL,
	authenticated_at TIMESTAMP NULL,
	subject VARCHAR(255) NOT NULL,
	remember BOOL NOT NULL DEFAULT false,
	nid UUID NOT NULL,
	identity_provider_session_id VARCHAR(40) NULL,
	expires_at TIMESTAMP NULL,
	CONSTRAINT hydra_oauth2_authentication_session_pkey PRIMARY KEY (id ASC),
	INDEX hydra_oauth2_authentication_session_subject_idx (subject ASC, nid ASC)
);
CREATE TABLE public.hydra_oauth2_obfuscated_authentication_session (
	subject VARCHAR(255) NOT NULL,
	client_id VARCHAR(255) NOT NULL,
	subject_obfuscated VARCHAR(255) NOT NULL,
	nid UUID NOT NULL,
	CONSTRAINT hydra_oauth2_obfuscated_authentication_session_pkey PRIMARY KEY (subject ASC, client_id ASC, nid ASC),
	UNIQUE INDEX hydra_oauth2_obfuscated_authentication_session_client_id_subject_obfuscated_idx (client_id ASC, subject_obfuscated ASC, nid ASC)
);
CREATE TABLE public.hydra_oauth2_logout_request (
	challenge VARCHAR(36) NOT NULL,
	verifier VARCHAR(36) NOT NULL,
	subject VARCHAR(255) NOT NULL,
	sid VARCHAR(36) NOT NULL,
	client_id VARCHAR(255) NULL,
	request_url STRING NOT NULL,
	redir_url STRING NOT NULL,
	was_used BOOL NOT NULL DEFAULT false,
	accepted BOOL NOT NULL DEFAULT false,
	rejected BOOL NOT NULL DEFAULT false,
	rp_initiated BOOL NOT NULL DEFAULT false,
	nid UUID NOT NULL,
	expires_at TIMESTAMP NULL,
	requested_at TIMESTAMP NULL,
	CONSTRAINT hydra_oauth2_logout_request_pkey PRIMARY KEY (challenge ASC),
	UNIQUE INDEX hydra_oauth2_logout_request_verifier_key (verifier ASC),
	INDEX hydra_oauth2_logout_request_client_id_idx (client_id ASC, nid ASC)
);
CREATE TABLE public.hydra_oauth2_flow (
	login_challenge VARCHAR(40) NOT NULL,
	login_verifier VARCHAR(40) NOT NULL,
	login_csrf VARCHAR(40) NOT NULL,
	subject VARCHAR(255) NOT NULL,
	request_url STRING NOT NULL,
	login_skip BOOL NOT NULL,
	client_id VARCHAR(255) NOT NULL,
	requested_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	login_initialized_at TIMESTAMP NULL,
	oidc_context JSONB NOT NULL DEFAULT '{}':::JSONB,
	login_session_id VARCHAR(40) NULL,
	state INT8 NOT NULL,
	login_remember BOOL NOT NULL DEFAULT false,
	login_remember_for INT8 NOT NULL,
	login_error STRING NULL,
	acr STRING NOT NULL DEFAULT '':::STRING,
	login_authenticated_at TIMESTAMP NULL,
	login_was_used BOOL NOT NULL DEFAULT false,
	forced_subject_identifier VARCHAR(255) NOT NULL DEFAULT '':::STRING::VARCHAR,
	context JSONB NULL DEFAULT '{}':::JSONB,
	consent_challenge_id VARCHAR(40) NULL,
	consent_skip BOOL NOT NULL DEFAULT false,
	consent_verifier VARCHAR(40) NULL,
	consent_csrf VARCHAR(40) NULL,
	consent_remember BOOL NOT NULL DEFAULT false,
	consent_remember_for INT8 NULL,
	consent_handled_at TIMESTAMP NULL,
	consent_error STRING NULL,
	session_access_token JSONB NOT NULL DEFAULT '{}':::JSONB,
	session_id_token JSONB NOT NULL DEFAULT '{}':::JSONB,
	consent_was_used BOOL NOT NULL DEFAULT false,
	nid UUID NOT NULL,
	requested_scope JSONB NOT NULL,
	requested_at_audience JSONB NULL DEFAULT '[]':::JSONB,
	amr JSONB NULL DEFAULT '[]':::JSONB,
	granted_scope JSONB NULL,
	granted_at_audience JSONB NULL DEFAULT '[]':::JSONB,
	login_extend_session_lifespan BOOL NOT NULL DEFAULT false,
	identity_provider_session_id VARCHAR(40) NULL,
	device_challenge_id VARCHAR(255) NULL,
	device_code_request_id VARCHAR(255) NULL,
	device_verifier VARCHAR(40) NULL,
	device_csrf VARCHAR(40) NULL,
	device_was_used BOOL NULL,
	device_handled_at TIMESTAMP NULL,
	device_error VARCHAR(2048) NULL,
	expires_at TIMESTAMP NULL AS (IF(consent_remember_for > 0:::INT8, requested_at + ('00:00:01':::INTERVAL * consent_remember_for), NULL)) VIRTUAL,
	CONSTRAINT hydra_oauth2_flow_pkey PRIMARY KEY (login_challenge ASC),
	UNIQUE INDEX hydra_oauth2_flow_consent_challenge_idx (consent_challenge_id ASC),
	INDEX hydra_oauth2_flow_client_id_subject_idx (client_id ASC, nid ASC, subject ASC),
	INDEX hydra_oauth2_flow_cid_idx (client_id ASC, nid ASC),
	INDEX hydra_oauth2_flow_login_session_id_idx (login_session_id ASC, nid ASC),
	INDEX hydra_oauth2_flow_sub_idx (subject ASC, nid ASC),
	INDEX hydra_oauth2_flow_previous_consents_idx (subject ASC, client_id ASC, nid ASC, consent_skip ASC, consent_error ASC, consent_remember ASC),
	UNIQUE INDEX hydra_oauth2_flow_device_challenge_idx (device_challenge_id ASC),
	CONSTRAINT check_state_state_state_state_login_remember_login_remember_for_login_error_acr_login_was_used_context_amr_state_login_remember_login_remember_for_login_error_acr_login_was_used_context_amr_state_login_remember_login_remember_for_login_error_acr_login_was_used_context_amr_consent_challenge_id_consent_verifier_consent_skip_consent_csrf_state_login_remember_login_remember_for_login_error_acr_login_was_used_context_amr_consent_challenge_id_consent_verifier_consent_skip_consent_csrf_state_login_remember_login_remember_for_login_error_acr_login_was_used_context_amr_consent_challenge_id_consent_verifier_consent_skip_consent_csrf_granted_scope_consent_remember_consent_remember_for_consent_error_session_access_token_session_id_token_consent_was_used CHECK ((((((((state = 128:::INT8) OR (state = 129:::INT8)) OR (state = 1:::INT8)) OR ((state = 2:::INT8) AND (((((((login_remember IS NOT NULL) AND (login_remember_for IS NOT NULL)) AND (login_error IS NOT NULL)) AND (acr IS NOT NULL)) AND (login_was_used IS NOT NULL)) AND (context IS NOT NULL)) AND (amr IS NOT NULL)))) OR ((state = 3:::INT8) AND (((((((login_remember IS NOT NULL) AND (login_remember_for IS NOT NULL)) AND (login_error IS NOT NULL)) AND (acr IS NOT NULL)) AND (login_was_used IS NOT NULL)) AND (context IS NOT NULL)) AND (amr IS NOT NULL)))) OR ((state = 4:::INT8) AND (((((((((((login_remember IS NOT NULL) AND (login_remember_for IS NOT NULL)) AND (login_error IS NOT NULL)) AND (acr IS NOT NULL)) AND (login_was_used IS NOT NULL)) AND (context IS NOT NULL)) AND (amr IS NOT NULL)) AND (consent_challenge_id IS NOT NULL)) AND (consent_verifier IS NOT NULL)) AND (consent_skip IS NOT NULL)) AND (consent_csrf IS NOT NULL)))) OR ((state = 5:::INT8) AND (((((((((((login_remember IS NOT NULL) AND (login_remember_for IS NOT NULL)) AND (login_error IS NOT NULL)) AND (acr IS NOT NULL)) AND (login_was_used IS NOT NULL)) AND (context IS NOT NULL)) AND (amr IS NOT NULL)) AND (consent_challenge_id IS NOT NULL)) AND (consent_verifier IS NOT NULL)) AND (consent_skip IS NOT NULL)) AND (consent_csrf IS NOT NULL)))) OR ((state = 6:::INT8) AND ((((((((((((((((((login_remember IS NOT NULL) AND (login_remember_for IS NOT NULL)) AND (login_error IS NOT NULL)) AND (acr IS NOT NULL)) AND (login_was_used IS NOT NULL)) AND (context IS NOT NULL)) AND (amr IS NOT NULL)) AND (consent_challenge_id IS NOT NULL)) AND (consent_verifier IS NOT NULL)) AND (consent_skip IS NOT NULL)) AND (consent_csrf IS NOT NULL)) AND (granted_scope IS NOT NULL)) AND (consent_remember IS NOT NULL)) AND (consent_remember_for IS NOT NULL)) AND (consent_error IS NOT NULL)) AND (session_access_token IS NOT NULL)) AND (session_id_token IS NOT NULL)) AND (consent_was_used IS NOT NULL))))
);
CREATE TABLE public.hydra_oauth2_access (
	signature VARCHAR(255) NOT NULL,
	request_id VARCHAR(40) NOT NULL,
	requested_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	client_id VARCHAR(255) NOT NULL,
	scope STRING NOT NULL,
	granted_scope STRING NOT NULL,
	form_data STRING NOT NULL,
	session_data STRING NOT NULL,
	subject VARCHAR(255) NOT NULL DEFAULT '':::STRING,
	active BOOL NOT NULL DEFAULT true,
	requested_audience STRING NULL DEFAULT '':::STRING,
	granted_audience STRING NULL DEFAULT '':::STRING,
	challenge_id VARCHAR(40) NULL,
	nid UUID NOT NULL,
	expires_at TIMESTAMP NULL,
	CONSTRAINT "primary" PRIMARY KEY (signature ASC),
	INDEX hydra_oauth2_access_requested_at_idx (requested_at ASC, nid ASC),
	INDEX hydra_oauth2_access_client_id_idx (client_id ASC, nid ASC),
	INDEX hydra_oauth2_access_challenge_id_idx (challenge_id ASC),
	INDEX hydra_oauth2_access_request_id_idx (request_id ASC, nid ASC)
);
CREATE TABLE public.hydra_oauth2_refresh (
	signature VARCHAR(255) NOT NULL,
	request_id VARCHAR(40) NOT NULL,
	requested_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	client_id VARCHAR(255) NOT NULL,
	scope STRING NOT NULL,
	granted_scope STRING NOT NULL,
	form_data STRING NOT NULL,
	session_data STRING NOT NULL,
	subject VARCHAR(255) NOT NULL DEFAULT '':::STRING,
	active BOOL NOT NULL DEFAULT true,
	requested_audience STRING NULL DEFAULT '':::STRING,
	granted_audience STRING NULL DEFAULT '':::STRING,
	challenge_id VARCHAR(40) NULL,
	nid UUID NOT NULL,
	expires_at TIMESTAMP NULL,
	first_used_at TIMESTAMP NULL,
	access_token_signature VARCHAR(255) NULL,
	used_times INT4 NULL,
	CONSTRAINT "primary" PRIMARY KEY (signature ASC),
	INDEX hydra_oauth2_refresh_client_id_idx (client_id ASC, nid ASC),
	INDEX hydra_oauth2_refresh_challenge_id_idx (challenge_id ASC),
	INDEX hydra_oauth2_refresh_request_id_idx (request_id ASC),
	INDEX hydra_oauth2_refresh_requested_at_idx (nid ASC, requested_at ASC)
);
CREATE TABLE public.hydra_oauth2_code (
	signature VARCHAR(255) NOT NULL,
	request_id VARCHAR(40) NOT NULL,
	requested_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	client_id VARCHAR(255) NOT NULL,
	scope STRING NOT NULL,
	granted_scope STRING NOT NULL,
	form_data STRING NOT NULL,
	session_data STRING NOT NULL,
	subject VARCHAR(255) NOT NULL DEFAULT '':::STRING,
	active BOOL NOT NULL DEFAULT true,
	requested_audience STRING NULL DEFAULT '':::STRING,
	granted_audience STRING NULL DEFAULT '':::STRING,
	challenge_id VARCHAR(40) NULL,
	nid UUID NOT NULL,
	expires_at TIMESTAMP NULL,
	CONSTRAINT "primary" PRIMARY KEY (signature ASC),
	INDEX hydra_oauth2_code_client_id_idx (client_id ASC, nid ASC),
	INDEX hydra_oauth2_code_challenge_id_idx (challenge_id ASC, nid ASC)
);
CREATE TABLE public.hydra_oauth2_oidc (
	signature VARCHAR(255) NOT NULL,
	request_id VARCHAR(40) NOT NULL,
	requested_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	client_id VARCHAR(255) NOT NULL,
	scope STRING NOT NULL,
	granted_scope STRING NOT NULL,
	form_data STRING NOT NULL,
	session_data STRING NOT NULL,
	subject VARCHAR(255) NOT NULL DEFAULT '':::STRING,
	active BOOL NOT NULL DEFAULT true,
	requested_audience STRING NULL DEFAULT '':::STRING,
	granted_audience STRING NULL DEFAULT '':::STRING,
	challenge_id VARCHAR(40) NULL,
	nid UUID NOT NULL,
	expires_at TIMESTAMP NULL,
	CONSTRAINT "primary" PRIMARY KEY (signature ASC),
	INDEX hydra_oauth2_oidc_client_id_idx (client_id ASC, nid ASC),
	INDEX hydra_oauth2_oidc_challenge_id_idx (challenge_id ASC)
);
CREATE TABLE public.hydra_oauth2_pkce (
	signature VARCHAR(255) NOT NULL,
	request_id VARCHAR(40) NOT NULL,
	requested_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	client_id VARCHAR(255) NOT NULL,
	scope STRING NOT NULL,
	granted_scope STRING NOT NULL,
	form_data STRING NOT NULL,
	session_data STRING NOT NULL,
	subject VARCHAR(255) NOT NULL,
	active BOOL NOT NULL DEFAULT true,
	requested_audience STRING NULL DEFAULT '':::STRING,
	granted_audience STRING NULL DEFAULT '':::STRING,
	challenge_id VARCHAR(40) NULL,
	nid UUID NOT NULL,
	expires_at TIMESTAMP NULL,
	CONSTRAINT "primary" PRIMARY KEY (signature ASC),
	INDEX hydra_oauth2_pkce_client_id_idx (client_id ASC, nid ASC),
	INDEX hydra_oauth2_pkce_challenge_id_idx (challenge_id ASC)
);
CREATE TABLE public.hydra_oauth2_jti_blacklist (
	signature VARCHAR(64) NOT NULL,
	expires_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	nid UUID NOT NULL,
	CONSTRAINT hydra_oauth2_jti_blacklist_pkey PRIMARY KEY (signature ASC, nid ASC),
	INDEX hydra_oauth2_jti_blacklist_expires_at_idx (expires_at ASC, nid ASC)
);
CREATE TABLE public.hydra_oauth2_trusted_jwt_bearer_issuer (
	id UUID NOT NULL,
	issuer VARCHAR(255) NOT NULL,
	subject VARCHAR(255) NOT NULL,
	scope STRING NOT NULL,
	key_set VARCHAR(255) NOT NULL,
	key_id VARCHAR(255) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	expires_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	nid UUID NOT NULL,
	allow_any_subject BOOL NOT NULL DEFAULT false,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	INDEX hydra_oauth2_trusted_jwt_bearer_issuer_expires_at_idx (expires_at ASC),
	INDEX hydra_oauth2_trusted_jwt_bearer_issuer_nid_idx (id ASC, nid ASC),
	UNIQUE INDEX hydra_oauth2_trusted_jwt_bearer_issuer_nid_uq_idx (nid ASC, key_id ASC, issuer ASC, subject ASC)
);
CREATE TABLE public.hydra_oauth2_device_auth_codes (
	device_code_signature VARCHAR(255) NOT NULL,
	user_code_signature VARCHAR(255) NOT NULL,
	request_id VARCHAR(40) NOT NULL,
	requested_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
	client_id VARCHAR(255) NOT NULL,
	scope VARCHAR(1024) NOT NULL,
	granted_scope VARCHAR(1024) NOT NULL,
	form_data VARCHAR(4096) NOT NULL,
	session_data STRING NOT NULL,
	subject VARCHAR(255) NOT NULL DEFAULT '':::STRING,
	device_code_active BOOL NOT NULL DEFAULT true,
	user_code_state INT2 NOT NULL DEFAULT 0:::INT8,
	requested_audience VARCHAR(1024) NOT NULL,
	granted_audience VARCHAR(1024) NOT NULL,
	challenge_id VARCHAR(40) NULL,
	expires_at TIMESTAMP NULL,
	nid UUID NOT NULL,
	CONSTRAINT hydra_oauth2_device_auth_codes_pkey PRIMARY KEY (device_code_signature ASC, nid ASC),
	INDEX hydra_oauth2_device_auth_codes_request_id_idx (request_id ASC, nid ASC),
	INDEX hydra_oauth2_device_auth_codes_client_id_idx (client_id ASC, nid ASC),
	INDEX hydra_oauth2_device_auth_codes_challenge_id_idx (challenge_id ASC),
	UNIQUE INDEX hydra_oauth2_device_auth_codes_user_code_signature_idx (nid ASC, user_code_signature ASC)
);
ALTER TABLE public.hydra_client ADD CONSTRAINT hydra_client_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_jwk ADD CONSTRAINT hydra_jwk_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_authentication_session ADD CONSTRAINT hydra_oauth2_authentication_session_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_obfuscated_authentication_session ADD CONSTRAINT hydra_oauth2_obfuscated_authentication_session_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_obfuscated_authentication_session ADD CONSTRAINT hydra_oauth2_obfuscated_authentication_session_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES public.hydra_client(id, nid) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_logout_request ADD CONSTRAINT hydra_oauth2_logout_request_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_logout_request ADD CONSTRAINT hydra_oauth2_logout_request_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES public.hydra_client(id, nid) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_flow ADD CONSTRAINT hydra_oauth2_flow_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_flow ADD CONSTRAINT hydra_oauth2_flow_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES public.hydra_client(id, nid) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_flow ADD CONSTRAINT hydra_oauth2_flow_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES public.hydra_oauth2_authentication_session(id) ON DELETE SET NULL;
ALTER TABLE public.hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES public.hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES public.hydra_client(id, nid) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES public.hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES public.hydra_client(id, nid) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES public.hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES public.hydra_client(id, nid) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES public.hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES public.hydra_client(id, nid) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES public.hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_client_id_fk FOREIGN KEY (client_id, nid) REFERENCES public.hydra_client(id, nid) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_jti_blacklist ADD CONSTRAINT hydra_oauth2_jti_blacklist_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_trusted_jwt_bearer_issuer ADD CONSTRAINT hydra_oauth2_trusted_jwt_bearer_issuer_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_trusted_jwt_bearer_issuer ADD CONSTRAINT fk_key_set_ref_hydra_jwk FOREIGN KEY (key_set, key_id, nid) REFERENCES public.hydra_jwk(sid, kid, nid) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_device_auth_codes ADD CONSTRAINT hydra_oauth2_device_auth_codes_client_id_nid_fkey FOREIGN KEY (client_id, nid) REFERENCES public.hydra_client(id, nid) ON DELETE CASCADE;
ALTER TABLE public.hydra_oauth2_device_auth_codes ADD CONSTRAINT hydra_oauth2_device_auth_codes_nid_fkey FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE ON UPDATE RESTRICT;
ALTER TABLE public.hydra_oauth2_device_auth_codes ADD CONSTRAINT hydra_oauth2_device_auth_codes_challenge_id_fkey FOREIGN KEY (challenge_id) REFERENCES public.hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;
ALTER TABLE public.hydra_client VALIDATE CONSTRAINT hydra_client_nid_fk_idx;
ALTER TABLE public.hydra_jwk VALIDATE CONSTRAINT hydra_jwk_nid_fk_idx;
ALTER TABLE public.hydra_oauth2_authentication_session VALIDATE CONSTRAINT hydra_oauth2_authentication_session_nid_fk_idx;
ALTER TABLE public.hydra_oauth2_obfuscated_authentication_session VALIDATE CONSTRAINT hydra_oauth2_obfuscated_authentication_session_nid_fk_idx;
ALTER TABLE public.hydra_oauth2_obfuscated_authentication_session VALIDATE CONSTRAINT hydra_oauth2_obfuscated_authentication_session_client_id_fk;
ALTER TABLE public.hydra_oauth2_logout_request VALIDATE CONSTRAINT hydra_oauth2_logout_request_nid_fk_idx;
ALTER TABLE public.hydra_oauth2_logout_request VALIDATE CONSTRAINT hydra_oauth2_logout_request_client_id_fk;
ALTER TABLE public.hydra_oauth2_flow VALIDATE CONSTRAINT hydra_oauth2_flow_nid_fk_idx;
ALTER TABLE public.hydra_oauth2_flow VALIDATE CONSTRAINT hydra_oauth2_flow_client_id_fk;
ALTER TABLE public.hydra_oauth2_flow VALIDATE CONSTRAINT hydra_oauth2_flow_login_session_id_fk;
ALTER TABLE public.hydra_oauth2_access VALIDATE CONSTRAINT hydra_oauth2_access_challenge_id_fk;
ALTER TABLE public.hydra_oauth2_access VALIDATE CONSTRAINT hydra_oauth2_access_nid_fk_idx;
ALTER TABLE public.hydra_oauth2_access VALIDATE CONSTRAINT hydra_oauth2_access_client_id_fk;
ALTER TABLE public.hydra_oauth2_refresh VALIDATE CONSTRAINT hydra_oauth2_refresh_challenge_id_fk;
ALTER TABLE public.hydra_oauth2_refresh VALIDATE CONSTRAINT hydra_oauth2_refresh_nid_fk_idx;
ALTER TABLE public.hydra_oauth2_refresh VALIDATE CONSTRAINT hydra_oauth2_refresh_client_id_fk;
ALTER TABLE public.hydra_oauth2_code VALIDATE CONSTRAINT hydra_oauth2_code_challenge_id_fk;
ALTER TABLE public.hydra_oauth2_code VALIDATE CONSTRAINT hydra_oauth2_code_nid_fk_idx;
ALTER TABLE public.hydra_oauth2_code VALIDATE CONSTRAINT hydra_oauth2_code_client_id_fk;
ALTER TABLE public.hydra_oauth2_oidc VALIDATE CONSTRAINT hydra_oauth2_oidc_challenge_id_fk;
ALTER TABLE public.hydra_oauth2_oidc VALIDATE CONSTRAINT hydra_oauth2_oidc_nid_fk_idx;
ALTER TABLE public.hydra_oauth2_oidc VALIDATE CONSTRAINT hydra_oauth2_oidc_client_id_fk;
ALTER TABLE public.hydra_oauth2_pkce VALIDATE CONSTRAINT hydra_oauth2_pkce_challenge_id_fk;
ALTER TABLE public.hydra_oauth2_pkce VALIDATE CONSTRAINT hydra_oauth2_pkce_nid_fk_idx;
ALTER TABLE public.hydra_oauth2_pkce VALIDATE CONSTRAINT hydra_oauth2_pkce_client_id_fk;
ALTER TABLE public.hydra_oauth2_jti_blacklist VALIDATE CONSTRAINT hydra_oauth2_jti_blacklist_nid_fk_idx;
ALTER TABLE public.hydra_oauth2_trusted_jwt_bearer_issuer VALIDATE CONSTRAINT hydra_oauth2_trusted_jwt_bearer_issuer_nid_fk_idx;
ALTER TABLE public.hydra_oauth2_trusted_jwt_bearer_issuer VALIDATE CONSTRAINT fk_key_set_ref_hydra_jwk;
ALTER TABLE public.hydra_oauth2_device_auth_codes VALIDATE CONSTRAINT hydra_oauth2_device_auth_codes_client_id_nid_fkey;
ALTER TABLE public.hydra_oauth2_device_auth_codes VALIDATE CONSTRAINT hydra_oauth2_device_auth_codes_nid_fkey;
ALTER TABLE public.hydra_oauth2_device_auth_codes VALIDATE CONSTRAINT hydra_oauth2_device_auth_codes_challenge_id_fkey;

