ALTER TABLE hydra_oauth2_consent_request ADD requested_at_audience TEXT NULL;
ALTER TABLE hydra_oauth2_authentication_request ADD requested_at_audience TEXT NULL;
ALTER TABLE hydra_oauth2_consent_request_handled ADD granted_at_audience TEXT NULL;

UPDATE hydra_oauth2_consent_request SET requested_at_audience='';
UPDATE hydra_oauth2_authentication_request SET requested_at_audience='';
UPDATE hydra_oauth2_consent_request_handled SET granted_at_audience='';

ALTER TABLE hydra_oauth2_consent_request MODIFY requested_at_audience TEXT NOT NULL;
ALTER TABLE hydra_oauth2_authentication_request MODIFY requested_at_audience TEXT NOT NULL;
ALTER TABLE hydra_oauth2_consent_request_handled MODIFY granted_at_audience TEXT NOT NULL;

