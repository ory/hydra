CREATE TABLE IF NOT EXISTS hydra_oauth2_device_link_request
(
    challenge character varying(40) COLLATE pg_catalog."default" NOT NULL,
    verifier character varying(40) COLLATE pg_catalog."default" NOT NULL,
    device_code character varying(255) COLLATE pg_catalog."default" NOT NULL,
    requested_scope text COLLATE pg_catalog."default" NOT NULL,
    requested_at_audience text COLLATE pg_catalog."default" DEFAULT ''::text,
    oidc_context text COLLATE pg_catalog."default" NOT NULL,
    client_id character varying(255) COLLATE pg_catalog."default" NOT NULL,
    login_challenge character varying(40) COLLATE pg_catalog."default",
    requested_at timestamp without time zone NOT NULL DEFAULT now(),
    request_url text COLLATE pg_catalog."default" NOT NULL,
    
    CONSTRAINT hydra_oauth2_device_link_request_pkey PRIMARY KEY (challenge),
    CONSTRAINT hydra_oauth2_device_link_request_client_id_fk FOREIGN KEY (client_id)
        REFERENCES public.hydra_client (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT hydra_oauth2_device_link_request_login_challenge_fk FOREIGN KEY (login_challenge)
        REFERENCES public.hydra_oauth2_authentication_request (challenge) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL
);