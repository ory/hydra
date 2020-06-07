--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.17
-- Dumped by pg_dump version 11.8 (Ubuntu 11.8-1.pgdg18.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: hydra_client; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_client (
    id character varying(255) NOT NULL,
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
    client_secret_expires_at integer DEFAULT 0 NOT NULL,
    sector_identifier_uri text NOT NULL,
    jwks text NOT NULL,
    jwks_uri text NOT NULL,
    request_uris text NOT NULL,
    token_endpoint_auth_method character varying(25) DEFAULT ''::character varying NOT NULL,
    request_object_signing_alg character varying(10) DEFAULT ''::character varying NOT NULL,
    userinfo_signed_response_alg character varying(10) DEFAULT ''::character varying NOT NULL,
    subject_type character varying(15) DEFAULT ''::character varying NOT NULL,
    allowed_cors_origins text NOT NULL,
    pk integer NOT NULL,
    audience text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    frontchannel_logout_uri text DEFAULT ''::text NOT NULL,
    frontchannel_logout_session_required boolean DEFAULT false NOT NULL,
    post_logout_redirect_uris text DEFAULT ''::text NOT NULL,
    backchannel_logout_uri text DEFAULT ''::text NOT NULL,
    backchannel_logout_session_required boolean DEFAULT false NOT NULL,
    metadata text NOT NULL,
    token_endpoint_auth_signing_alg character varying(10) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.hydra_client OWNER TO postgres;

--
-- Name: hydra_client_pk_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.hydra_client_pk_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.hydra_client_pk_seq OWNER TO postgres;

--
-- Name: hydra_client_pk_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.hydra_client_pk_seq OWNED BY public.hydra_client.pk;


--
-- Name: hydra_jwk; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_jwk (
    sid character varying(255) NOT NULL,
    kid character varying(255) NOT NULL,
    version integer DEFAULT 0 NOT NULL,
    keydata text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    pk integer NOT NULL
);


ALTER TABLE public.hydra_jwk OWNER TO postgres;

--
-- Name: hydra_jwk_pk_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.hydra_jwk_pk_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.hydra_jwk_pk_seq OWNER TO postgres;

--
-- Name: hydra_jwk_pk_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.hydra_jwk_pk_seq OWNED BY public.hydra_jwk.pk;


--
-- Name: hydra_oauth2_access; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_access (
    signature character varying(255) NOT NULL,
    request_id character varying(40) NOT NULL,
    requested_at timestamp without time zone DEFAULT now() NOT NULL,
    client_id character varying(255) NOT NULL,
    scope text NOT NULL,
    granted_scope text NOT NULL,
    form_data text NOT NULL,
    session_data text NOT NULL,
    subject character varying(255) DEFAULT ''::character varying NOT NULL,
    active boolean DEFAULT true NOT NULL,
    requested_audience text DEFAULT ''::text,
    granted_audience text DEFAULT ''::text,
    challenge_id character varying(40)
);


ALTER TABLE public.hydra_oauth2_access OWNER TO postgres;

--
-- Name: hydra_oauth2_authentication_request; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_authentication_request (
    challenge character varying(40) NOT NULL,
    requested_scope text NOT NULL,
    verifier character varying(40) NOT NULL,
    csrf character varying(40) NOT NULL,
    subject character varying(255) NOT NULL,
    request_url text NOT NULL,
    skip boolean NOT NULL,
    client_id character varying(255) NOT NULL,
    requested_at timestamp without time zone DEFAULT now() NOT NULL,
    authenticated_at timestamp without time zone,
    oidc_context text NOT NULL,
    login_session_id character varying(40),
    requested_at_audience text DEFAULT ''::text
);


ALTER TABLE public.hydra_oauth2_authentication_request OWNER TO postgres;

--
-- Name: hydra_oauth2_authentication_request_handled; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_authentication_request_handled (
    challenge character varying(40) NOT NULL,
    subject character varying(255) NOT NULL,
    remember boolean NOT NULL,
    remember_for integer NOT NULL,
    error text NOT NULL,
    acr text NOT NULL,
    requested_at timestamp without time zone DEFAULT now() NOT NULL,
    authenticated_at timestamp without time zone,
    was_used boolean NOT NULL,
    forced_subject_identifier character varying(255) DEFAULT ''::character varying,
    context text DEFAULT '{}'::text NOT NULL
);


ALTER TABLE public.hydra_oauth2_authentication_request_handled OWNER TO postgres;

--
-- Name: hydra_oauth2_authentication_session; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_authentication_session (
    id character varying(40) NOT NULL,
    authenticated_at timestamp without time zone DEFAULT now() NOT NULL,
    subject character varying(255) NOT NULL,
    remember boolean DEFAULT false NOT NULL
);


ALTER TABLE public.hydra_oauth2_authentication_session OWNER TO postgres;

--
-- Name: hydra_oauth2_code; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_code (
    signature character varying(255) NOT NULL,
    request_id character varying(40) NOT NULL,
    requested_at timestamp without time zone DEFAULT now() NOT NULL,
    client_id character varying(255) NOT NULL,
    scope text NOT NULL,
    granted_scope text NOT NULL,
    form_data text NOT NULL,
    session_data text NOT NULL,
    subject character varying(255) DEFAULT ''::character varying NOT NULL,
    active boolean DEFAULT true NOT NULL,
    requested_audience text DEFAULT ''::text,
    granted_audience text DEFAULT ''::text,
    challenge_id character varying(40)
);


ALTER TABLE public.hydra_oauth2_code OWNER TO postgres;

--
-- Name: hydra_oauth2_consent_request; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_consent_request (
    challenge character varying(40) NOT NULL,
    verifier character varying(40) NOT NULL,
    client_id character varying(255) NOT NULL,
    subject character varying(255) NOT NULL,
    request_url text NOT NULL,
    skip boolean NOT NULL,
    requested_scope text NOT NULL,
    csrf character varying(40) NOT NULL,
    authenticated_at timestamp without time zone,
    requested_at timestamp without time zone DEFAULT now() NOT NULL,
    oidc_context text NOT NULL,
    forced_subject_identifier character varying(255) DEFAULT ''::character varying,
    login_session_id character varying(40),
    login_challenge character varying(40),
    requested_at_audience text DEFAULT ''::text,
    acr text DEFAULT ''::text,
    context text DEFAULT '{}'::text NOT NULL
);


ALTER TABLE public.hydra_oauth2_consent_request OWNER TO postgres;

--
-- Name: hydra_oauth2_consent_request_handled; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_consent_request_handled (
    challenge character varying(40) NOT NULL,
    granted_scope text NOT NULL,
    remember boolean NOT NULL,
    remember_for integer NOT NULL,
    error text NOT NULL,
    requested_at timestamp without time zone DEFAULT now() NOT NULL,
    session_access_token text NOT NULL,
    session_id_token text NOT NULL,
    authenticated_at timestamp without time zone,
    was_used boolean NOT NULL,
    granted_at_audience text DEFAULT ''::text,
    handled_at timestamp without time zone
);


ALTER TABLE public.hydra_oauth2_consent_request_handled OWNER TO postgres;

--
-- Name: hydra_oauth2_jti_blacklist; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_jti_blacklist (
    signature character varying(64) NOT NULL,
    expires_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.hydra_oauth2_jti_blacklist OWNER TO postgres;

--
-- Name: hydra_oauth2_logout_request; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_logout_request (
    challenge character varying(36) NOT NULL,
    verifier character varying(36) NOT NULL,
    subject character varying(255) NOT NULL,
    sid character varying(36) NOT NULL,
    client_id character varying(255),
    request_url text NOT NULL,
    redir_url text NOT NULL,
    was_used boolean DEFAULT false NOT NULL,
    accepted boolean DEFAULT false NOT NULL,
    rejected boolean DEFAULT false NOT NULL,
    rp_initiated boolean DEFAULT false NOT NULL
);


ALTER TABLE public.hydra_oauth2_logout_request OWNER TO postgres;

--
-- Name: hydra_oauth2_obfuscated_authentication_session; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_obfuscated_authentication_session (
    subject character varying(255) NOT NULL,
    client_id character varying(255) NOT NULL,
    subject_obfuscated character varying(255) NOT NULL
);


ALTER TABLE public.hydra_oauth2_obfuscated_authentication_session OWNER TO postgres;

--
-- Name: hydra_oauth2_oidc; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_oidc (
    signature character varying(255) NOT NULL,
    request_id character varying(40) NOT NULL,
    requested_at timestamp without time zone DEFAULT now() NOT NULL,
    client_id character varying(255) NOT NULL,
    scope text NOT NULL,
    granted_scope text NOT NULL,
    form_data text NOT NULL,
    session_data text NOT NULL,
    subject character varying(255) DEFAULT ''::character varying NOT NULL,
    active boolean DEFAULT true NOT NULL,
    requested_audience text DEFAULT ''::text,
    granted_audience text DEFAULT ''::text,
    challenge_id character varying(40)
);


ALTER TABLE public.hydra_oauth2_oidc OWNER TO postgres;

--
-- Name: hydra_oauth2_pkce; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_pkce (
    signature character varying(255) NOT NULL,
    request_id character varying(40) NOT NULL,
    requested_at timestamp without time zone DEFAULT now() NOT NULL,
    client_id character varying(255) NOT NULL,
    scope text NOT NULL,
    granted_scope text NOT NULL,
    form_data text NOT NULL,
    session_data text NOT NULL,
    subject character varying(255) NOT NULL,
    active boolean DEFAULT true NOT NULL,
    requested_audience text DEFAULT ''::text,
    granted_audience text DEFAULT ''::text,
    challenge_id character varying(40)
);


ALTER TABLE public.hydra_oauth2_pkce OWNER TO postgres;

--
-- Name: hydra_oauth2_refresh; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hydra_oauth2_refresh (
    signature character varying(255) NOT NULL,
    request_id character varying(40) NOT NULL,
    requested_at timestamp without time zone DEFAULT now() NOT NULL,
    client_id character varying(255) NOT NULL,
    scope text NOT NULL,
    granted_scope text NOT NULL,
    form_data text NOT NULL,
    session_data text NOT NULL,
    subject character varying(255) DEFAULT ''::character varying NOT NULL,
    active boolean DEFAULT true NOT NULL,
    requested_audience text DEFAULT ''::text,
    granted_audience text DEFAULT ''::text,
    challenge_id character varying(40)
);


ALTER TABLE public.hydra_oauth2_refresh OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: hydra_client pk; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_client ALTER COLUMN pk SET DEFAULT nextval('public.hydra_client_pk_seq'::regclass);


--
-- Name: hydra_jwk pk; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_jwk ALTER COLUMN pk SET DEFAULT nextval('public.hydra_jwk_pk_seq'::regclass);


--
-- Name: hydra_client hydra_client_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_client
    ADD CONSTRAINT hydra_client_pkey PRIMARY KEY (pk);


--
-- Name: hydra_jwk hydra_jwk_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_jwk
    ADD CONSTRAINT hydra_jwk_pkey PRIMARY KEY (pk);


--
-- Name: hydra_oauth2_access hydra_oauth2_access_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_access
    ADD CONSTRAINT hydra_oauth2_access_pkey PRIMARY KEY (signature);


--
-- Name: hydra_oauth2_authentication_request_handled hydra_oauth2_authentication_request_handled_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_authentication_request_handled
    ADD CONSTRAINT hydra_oauth2_authentication_request_handled_pkey PRIMARY KEY (challenge);


--
-- Name: hydra_oauth2_authentication_request hydra_oauth2_authentication_request_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_authentication_request
    ADD CONSTRAINT hydra_oauth2_authentication_request_pkey PRIMARY KEY (challenge);


--
-- Name: hydra_oauth2_authentication_session hydra_oauth2_authentication_session_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_authentication_session
    ADD CONSTRAINT hydra_oauth2_authentication_session_pkey PRIMARY KEY (id);


--
-- Name: hydra_oauth2_code hydra_oauth2_code_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_code
    ADD CONSTRAINT hydra_oauth2_code_pkey PRIMARY KEY (signature);


--
-- Name: hydra_oauth2_consent_request_handled hydra_oauth2_consent_request_handled_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_consent_request_handled
    ADD CONSTRAINT hydra_oauth2_consent_request_handled_pkey PRIMARY KEY (challenge);


--
-- Name: hydra_oauth2_consent_request hydra_oauth2_consent_request_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_consent_request
    ADD CONSTRAINT hydra_oauth2_consent_request_pkey PRIMARY KEY (challenge);


--
-- Name: hydra_oauth2_jti_blacklist hydra_oauth2_jti_blacklist_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_jti_blacklist
    ADD CONSTRAINT hydra_oauth2_jti_blacklist_pkey PRIMARY KEY (signature);


--
-- Name: hydra_oauth2_logout_request hydra_oauth2_logout_request_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_logout_request
    ADD CONSTRAINT hydra_oauth2_logout_request_pkey PRIMARY KEY (challenge);


--
-- Name: hydra_oauth2_obfuscated_authentication_session hydra_oauth2_obfuscated_authentication_session_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_obfuscated_authentication_session
    ADD CONSTRAINT hydra_oauth2_obfuscated_authentication_session_pkey PRIMARY KEY (subject, client_id);


--
-- Name: hydra_oauth2_oidc hydra_oauth2_oidc_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_oidc
    ADD CONSTRAINT hydra_oauth2_oidc_pkey PRIMARY KEY (signature);


--
-- Name: hydra_oauth2_pkce hydra_oauth2_pkce_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_pkce
    ADD CONSTRAINT hydra_oauth2_pkce_pkey PRIMARY KEY (signature);


--
-- Name: hydra_oauth2_refresh hydra_oauth2_refresh_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_refresh
    ADD CONSTRAINT hydra_oauth2_refresh_pkey PRIMARY KEY (signature);


--
-- Name: hydra_client_idx_id_uq; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hydra_client_idx_id_uq ON public.hydra_client USING btree (id);


--
-- Name: hydra_jwk_idx_id_uq; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hydra_jwk_idx_id_uq ON public.hydra_jwk USING btree (sid, kid);


--
-- Name: hydra_oauth2_access_challenge_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_access_challenge_id_idx ON public.hydra_oauth2_access USING btree (challenge_id);


--
-- Name: hydra_oauth2_access_client_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_access_client_id_idx ON public.hydra_oauth2_access USING btree (client_id);


--
-- Name: hydra_oauth2_access_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hydra_oauth2_access_request_id_idx ON public.hydra_oauth2_access USING btree (request_id);


--
-- Name: hydra_oauth2_access_requested_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_access_requested_at_idx ON public.hydra_oauth2_access USING btree (requested_at);


--
-- Name: hydra_oauth2_authentication_request_cid_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_authentication_request_cid_idx ON public.hydra_oauth2_authentication_request USING btree (client_id);


--
-- Name: hydra_oauth2_authentication_request_login_session_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_authentication_request_login_session_id_idx ON public.hydra_oauth2_authentication_request USING btree (login_session_id);


--
-- Name: hydra_oauth2_authentication_request_sub_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_authentication_request_sub_idx ON public.hydra_oauth2_authentication_request USING btree (subject);


--
-- Name: hydra_oauth2_authentication_request_veri_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hydra_oauth2_authentication_request_veri_idx ON public.hydra_oauth2_authentication_request USING btree (verifier);


--
-- Name: hydra_oauth2_authentication_session_sub_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_authentication_session_sub_idx ON public.hydra_oauth2_authentication_session USING btree (subject);


--
-- Name: hydra_oauth2_code_challenge_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_code_challenge_id_idx ON public.hydra_oauth2_code USING btree (challenge_id);


--
-- Name: hydra_oauth2_code_client_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_code_client_id_idx ON public.hydra_oauth2_code USING btree (client_id);


--
-- Name: hydra_oauth2_consent_request_cid_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_consent_request_cid_idx ON public.hydra_oauth2_consent_request USING btree (client_id);


--
-- Name: hydra_oauth2_consent_request_client_id_subject_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_consent_request_client_id_subject_idx ON public.hydra_oauth2_consent_request USING btree (client_id, subject);


--
-- Name: hydra_oauth2_consent_request_login_challenge_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_consent_request_login_challenge_idx ON public.hydra_oauth2_consent_request USING btree (login_challenge);


--
-- Name: hydra_oauth2_consent_request_login_session_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_consent_request_login_session_id_idx ON public.hydra_oauth2_consent_request USING btree (login_session_id);


--
-- Name: hydra_oauth2_consent_request_sub_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_consent_request_sub_idx ON public.hydra_oauth2_consent_request USING btree (subject);


--
-- Name: hydra_oauth2_consent_request_veri_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hydra_oauth2_consent_request_veri_idx ON public.hydra_oauth2_consent_request USING btree (verifier);


--
-- Name: hydra_oauth2_jti_blacklist_expiry; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_jti_blacklist_expiry ON public.hydra_oauth2_jti_blacklist USING btree (expires_at);


--
-- Name: hydra_oauth2_logout_request_client_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_logout_request_client_id_idx ON public.hydra_oauth2_logout_request USING btree (client_id);


--
-- Name: hydra_oauth2_logout_request_veri_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hydra_oauth2_logout_request_veri_idx ON public.hydra_oauth2_logout_request USING btree (verifier);


--
-- Name: hydra_oauth2_obfuscated_authentication_session_so_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hydra_oauth2_obfuscated_authentication_session_so_idx ON public.hydra_oauth2_obfuscated_authentication_session USING btree (client_id, subject_obfuscated);


--
-- Name: hydra_oauth2_oidc_challenge_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_oidc_challenge_id_idx ON public.hydra_oauth2_oidc USING btree (challenge_id);


--
-- Name: hydra_oauth2_oidc_client_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_oidc_client_id_idx ON public.hydra_oauth2_oidc USING btree (client_id);


--
-- Name: hydra_oauth2_pkce_challenge_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_pkce_challenge_id_idx ON public.hydra_oauth2_pkce USING btree (challenge_id);


--
-- Name: hydra_oauth2_pkce_client_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_pkce_client_id_idx ON public.hydra_oauth2_pkce USING btree (client_id);


--
-- Name: hydra_oauth2_refresh_challenge_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_refresh_challenge_id_idx ON public.hydra_oauth2_refresh USING btree (challenge_id);


--
-- Name: hydra_oauth2_refresh_client_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hydra_oauth2_refresh_client_id_idx ON public.hydra_oauth2_refresh USING btree (client_id);


--
-- Name: hydra_oauth2_refresh_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hydra_oauth2_refresh_request_id_idx ON public.hydra_oauth2_refresh USING btree (request_id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: hydra_oauth2_access hydra_oauth2_access_challenge_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_access
    ADD CONSTRAINT hydra_oauth2_access_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES public.hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_access hydra_oauth2_access_client_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_access
    ADD CONSTRAINT hydra_oauth2_access_client_id_fk FOREIGN KEY (client_id) REFERENCES public.hydra_client(id) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_authentication_request hydra_oauth2_authentication_request_client_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_authentication_request
    ADD CONSTRAINT hydra_oauth2_authentication_request_client_id_fk FOREIGN KEY (client_id) REFERENCES public.hydra_client(id) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_authentication_request_handled hydra_oauth2_authentication_request_handled_challenge_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_authentication_request_handled
    ADD CONSTRAINT hydra_oauth2_authentication_request_handled_challenge_fk FOREIGN KEY (challenge) REFERENCES public.hydra_oauth2_authentication_request(challenge) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_authentication_request hydra_oauth2_authentication_request_login_session_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_authentication_request
    ADD CONSTRAINT hydra_oauth2_authentication_request_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES public.hydra_oauth2_authentication_session(id) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_code hydra_oauth2_code_challenge_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_code
    ADD CONSTRAINT hydra_oauth2_code_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES public.hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_code hydra_oauth2_code_client_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_code
    ADD CONSTRAINT hydra_oauth2_code_client_id_fk FOREIGN KEY (client_id) REFERENCES public.hydra_client(id) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_consent_request hydra_oauth2_consent_request_client_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_consent_request
    ADD CONSTRAINT hydra_oauth2_consent_request_client_id_fk FOREIGN KEY (client_id) REFERENCES public.hydra_client(id) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_consent_request_handled hydra_oauth2_consent_request_handled_challenge_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_consent_request_handled
    ADD CONSTRAINT hydra_oauth2_consent_request_handled_challenge_fk FOREIGN KEY (challenge) REFERENCES public.hydra_oauth2_consent_request(challenge) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_consent_request hydra_oauth2_consent_request_login_challenge_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_consent_request
    ADD CONSTRAINT hydra_oauth2_consent_request_login_challenge_fk FOREIGN KEY (login_challenge) REFERENCES public.hydra_oauth2_authentication_request(challenge) ON DELETE SET NULL;


--
-- Name: hydra_oauth2_consent_request hydra_oauth2_consent_request_login_session_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_consent_request
    ADD CONSTRAINT hydra_oauth2_consent_request_login_session_id_fk FOREIGN KEY (login_session_id) REFERENCES public.hydra_oauth2_authentication_session(id) ON DELETE SET NULL;


--
-- Name: hydra_oauth2_logout_request hydra_oauth2_logout_request_client_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_logout_request
    ADD CONSTRAINT hydra_oauth2_logout_request_client_id_fk FOREIGN KEY (client_id) REFERENCES public.hydra_client(id) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_obfuscated_authentication_session hydra_oauth2_obfuscated_authentication_session_client_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_obfuscated_authentication_session
    ADD CONSTRAINT hydra_oauth2_obfuscated_authentication_session_client_id_fk FOREIGN KEY (client_id) REFERENCES public.hydra_client(id) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_oidc hydra_oauth2_oidc_challenge_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_oidc
    ADD CONSTRAINT hydra_oauth2_oidc_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES public.hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_oidc hydra_oauth2_oidc_client_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_oidc
    ADD CONSTRAINT hydra_oauth2_oidc_client_id_fk FOREIGN KEY (client_id) REFERENCES public.hydra_client(id) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_pkce hydra_oauth2_pkce_challenge_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_pkce
    ADD CONSTRAINT hydra_oauth2_pkce_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES public.hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_pkce hydra_oauth2_pkce_client_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_pkce
    ADD CONSTRAINT hydra_oauth2_pkce_client_id_fk FOREIGN KEY (client_id) REFERENCES public.hydra_client(id) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_refresh hydra_oauth2_refresh_challenge_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_refresh
    ADD CONSTRAINT hydra_oauth2_refresh_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES public.hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;


--
-- Name: hydra_oauth2_refresh hydra_oauth2_refresh_client_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hydra_oauth2_refresh
    ADD CONSTRAINT hydra_oauth2_refresh_client_id_fk FOREIGN KEY (client_id) REFERENCES public.hydra_client(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

