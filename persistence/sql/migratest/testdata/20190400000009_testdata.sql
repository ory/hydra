INSERT INTO hydra_oauth2_access (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
SELECT 'sig-0009', 'req-0009', '2022-02-15 22:20:22', hc.id, 'scope-0009', 'granted_scope-0009', 'form_data-0009', 'session-0009', 'subject-0009', false, 'requested_audience-0009', 'granted_audience-0009', crh.challenge
FROM hydra_client hc, hydra_oauth2_consent_request_handled crh
ORDER BY hc.pk, crh.challenge DESC
LIMIT 1;

INSERT INTO hydra_oauth2_refresh (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
SELECT 'sig-0009', 'req-0009', '2022-02-15 22:20:22', hc.id, 'scope-0009', 'granted_scope-0009', 'form_data-0009', 'session-0009', 'subject-0009', false, 'requested_audience-0009', 'granted_audience-0009', crh.challenge
FROM hydra_client hc, hydra_oauth2_consent_request_handled crh
ORDER BY hc.pk, crh.challenge DESC
LIMIT 1;

INSERT INTO hydra_oauth2_code (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
SELECT 'sig-0009', 'req-0009', '2022-02-15 22:20:22', hc.id, 'scope-0009', 'granted_scope-0009', 'form_data-0009', 'session-0009', 'subject-0009', false, 'requested_audience-0009', 'granted_audience-0009', crh.challenge
FROM hydra_client hc, hydra_oauth2_consent_request_handled crh
ORDER BY hc.pk, crh.challenge DESC
LIMIT 1;

INSERT INTO hydra_oauth2_oidc (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
SELECT 'sig-0009', 'req-0009', '2022-02-15 22:20:22', hc.id, 'scope-0009', 'granted_scope-0009', 'form_data-0009', 'session-0009', 'subject-0009', false, 'requested_audience-0009', 'granted_audience-0009', crh.challenge
FROM hydra_client hc, hydra_oauth2_consent_request_handled crh
ORDER BY hc.pk, crh.challenge DESC
LIMIT 1;

INSERT INTO hydra_oauth2_pkce (signature, request_id, requested_at, client_id, scope, granted_scope, form_data, session_data, subject, active, requested_audience, granted_audience, challenge_id)
SELECT 'sig-0009', 'req-0009', '2022-02-15 22:20:22', hc.id, 'scope-0009', 'granted_scope-0009', 'form_data-0009', 'session-0009', 'subject-0009', false, 'requested_audience-0009', 'granted_audience-0009', crh.challenge
FROM hydra_client hc, hydra_oauth2_consent_request_handled crh
ORDER BY hc.pk, crh.challenge DESC
LIMIT 1;
