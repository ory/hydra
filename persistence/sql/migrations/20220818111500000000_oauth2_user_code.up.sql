CREATE TABLE IF NOT EXISTS hydra_oauth2_user_code (
    signature character varying(255) NOT NULL PRIMARY KEY,
    user_code character varying(40),
    device_code character varying(40),
    request_id character varying(40) NOT NULL,
    requested_at timestamp without time zone NOT NULL DEFAULT now(),
    client_id character varying(255) NOT NULL,
    scope text NOT NULL,
    granted_scope text NOT NULL,
    form_data text NOT NULL,
    session_data text NOT NULL,
    subject character varying(255)  NOT NULL,
    active boolean NOT NULL,
    requested_audience text,
    granted_audience text,
    challenge_id character varying(40) ,
    -- CONSTRAINT hydra_oauth2_code_challenge_id_fk FOREIGN KEY (challenge_id)
    --     REFERENCES public.hydra_oauth2_consent_request_handled (challenge) MATCH SIMPLE
    --     ON UPDATE NO ACTION
    --     ON DELETE CASCADE,
    CONSTRAINT hydra_oauth2_device_client_id_fk FOREIGN KEY (client_id)
        REFERENCES public.hydra_client (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)
