-- +migrate Up
CREATE UNIQUE INDEX hydra_oauth2_obfuscated_authentication_session_so_idx ON hydra_oauth2_obfuscated_authentication_session (client_id, subject_obfuscated);

CREATE INDEX hydra_oauth2_consent_request_cid_idx ON hydra_oauth2_consent_request (client_id);
CREATE INDEX hydra_oauth2_consent_request_sub_idx ON hydra_oauth2_consent_request (subject);
CREATE UNIQUE INDEX hydra_oauth2_consent_request_veri_idx ON hydra_oauth2_consent_request (verifier);

CREATE INDEX hydra_oauth2_authentication_request_cid_idx ON hydra_oauth2_authentication_request (client_id);
CREATE INDEX hydra_oauth2_authentication_request_sub_idx ON hydra_oauth2_authentication_request (subject);
CREATE UNIQUE INDEX 5hydra_oauth2_authentication_request_veri_idx ON hydra_oauth2_authentication_request (verifier);

-- +migrate Down
DROP INDEX hydra_oauth2_obfuscated_authentication_session_so_idx ON hydra_oauth2_obfuscated_authentication_session;

DROP INDEX hydra_oauth2_consent_request_cid_idx ON hydra_oauth2_consent_request;
DROP INDEX hydra_oauth2_consent_request_sub_idx ON hydra_oauth2_consent_request;
DROP INDEX hydra_oauth2_consent_request_veri_idx ON hydra_oauth2_consent_request;

DROP INDEX hydra_oauth2_authentication_request_cid_idx ON hydra_oauth2_authentication_request;
DROP INDEX hydra_oauth2_authentication_request_sub_idx ON hydra_oauth2_authentication_request;
DROP INDEX hydra_oauth2_authentication_request_veri_idx ON hydra_oauth2_authentication_request;
