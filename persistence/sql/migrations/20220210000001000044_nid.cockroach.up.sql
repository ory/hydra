-- Migration generated by the command below; DO NOT EDIT.
-- hydra:generate hydra migrate gen

ALTER TABLE hydra_oauth2_obfuscated_authentication_session DROP CONSTRAINT "primary";
ALTER TABLE hydra_oauth2_obfuscated_authentication_session ADD CONSTRAINT "hydra_oauth2_obfuscated_authentication_session_pkey" PRIMARY KEY (subject ASC, client_id ASC, nid ASC);
