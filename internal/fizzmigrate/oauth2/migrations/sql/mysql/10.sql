-- +migrate Up
CREATE INDEX hydra_oauth2_access_challenge_id_idx ON hydra_oauth2_access (challenge_id);
CREATE INDEX hydra_oauth2_refresh_challenge_id_idx ON hydra_oauth2_refresh (challenge_id);
CREATE INDEX hydra_oauth2_code_challenge_id_idx ON hydra_oauth2_code (challenge_id);
CREATE INDEX hydra_oauth2_oidc_challenge_id_idx ON hydra_oauth2_oidc (challenge_id);
CREATE INDEX hydra_oauth2_pkce_challenge_id_idx ON hydra_oauth2_pkce (challenge_id);

-- +migrate Down
ALTER TABLE hydra_oauth2_access DROP FOREIGN KEY hydra_oauth2_access_challenge_id_fk;
DROP INDEX hydra_oauth2_access_challenge_id_idx ON hydra_oauth2_access;
ALTER TABLE hydra_oauth2_access ADD CONSTRAINT hydra_oauth2_access_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;

ALTER TABLE hydra_oauth2_refresh DROP FOREIGN KEY hydra_oauth2_refresh_challenge_id_fk;
DROP INDEX hydra_oauth2_refresh_challenge_id_idx ON hydra_oauth2_refresh;
ALTER TABLE hydra_oauth2_refresh ADD CONSTRAINT hydra_oauth2_refresh_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;

ALTER TABLE hydra_oauth2_code DROP FOREIGN KEY hydra_oauth2_code_challenge_id_fk;
DROP INDEX hydra_oauth2_code_challenge_id_idx ON hydra_oauth2_code;
ALTER TABLE hydra_oauth2_code ADD CONSTRAINT hydra_oauth2_code_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;

ALTER TABLE hydra_oauth2_oidc DROP FOREIGN KEY hydra_oauth2_oidc_challenge_id_fk;
DROP INDEX hydra_oauth2_oidc_challenge_id_idx ON hydra_oauth2_oidc;
ALTER TABLE hydra_oauth2_oidc ADD CONSTRAINT hydra_oauth2_oidc_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;

ALTER TABLE hydra_oauth2_pkce DROP FOREIGN KEY hydra_oauth2_pkce_challenge_id_fk;
DROP INDEX hydra_oauth2_pkce_challenge_id_idx ON hydra_oauth2_pkce;
ALTER TABLE hydra_oauth2_pkce ADD CONSTRAINT hydra_oauth2_pkce_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_consent_request_handled(challenge) ON DELETE CASCADE;
