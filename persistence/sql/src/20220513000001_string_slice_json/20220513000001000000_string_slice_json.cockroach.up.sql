ALTER TABLE hydra_client ADD COLUMN redirect_uris_json jsonb DEFAULT '[]' NOT NULL;
ALTER TABLE hydra_client ADD COLUMN grant_types_json jsonb DEFAULT '[]' NOT NULL;
ALTER TABLE hydra_client ADD COLUMN response_types_json jsonb DEFAULT '[]' NOT NULL;
ALTER TABLE hydra_client ADD COLUMN audience_json jsonb DEFAULT '[]' NOT NULL;
ALTER TABLE hydra_client ADD COLUMN allowed_cors_origins_json jsonb DEFAULT '[]' NOT NULL;
ALTER TABLE hydra_client ADD COLUMN contacts_json jsonb DEFAULT '[]' NOT NULL;
ALTER TABLE hydra_client ADD COLUMN request_uris_json jsonb DEFAULT '[]' NOT NULL;
ALTER TABLE hydra_client ADD COLUMN post_logout_redirect_uris_json jsonb DEFAULT '[]' NOT NULL;
--split
UPDATE hydra_client SET redirect_uris_json = cast('["' || REPLACE(redirect_uris,'|','","') || '"]' as jsonb) WHERE redirect_uris <> '';
UPDATE hydra_client SET grant_types_json = cast('["' || REPLACE(grant_types,'|','","') || '"]' as jsonb) WHERE grant_types <> '';
UPDATE hydra_client SET response_types_json = cast('["' || REPLACE(response_types,'|','","') || '"]' as jsonb) WHERE response_types <> '';
UPDATE hydra_client SET audience_json = cast('["' || REPLACE(audience,'|','","') || '"]' as jsonb) WHERE audience <> '';
UPDATE hydra_client SET allowed_cors_origins_json = cast('["' || REPLACE(allowed_cors_origins,'|','","') || '"]' as jsonb) WHERE allowed_cors_origins <> '';
UPDATE hydra_client SET contacts_json = cast('["' || REPLACE(contacts,'|','","') || '"]' as jsonb) WHERE contacts <> '';
UPDATE hydra_client SET request_uris_json = cast('["' || REPLACE(request_uris,'|','","') || '"]' as jsonb) WHERE request_uris <> '';
UPDATE hydra_client SET post_logout_redirect_uris_json = cast('["' || REPLACE(post_logout_redirect_uris,'|','","') || '"]' as jsonb) WHERE post_logout_redirect_uris <> '';
--split
ALTER TABLE hydra_client ALTER COLUMN redirect_uris_json DROP DEFAULT;
ALTER TABLE hydra_client ALTER COLUMN grant_types_json DROP DEFAULT;
ALTER TABLE hydra_client ALTER COLUMN response_types_json DROP DEFAULT;
ALTER TABLE hydra_client ALTER COLUMN audience_json DROP DEFAULT;
ALTER TABLE hydra_client ALTER COLUMN allowed_cors_origins_json DROP DEFAULT;
ALTER TABLE hydra_client ALTER COLUMN contacts_json DROP DEFAULT;
ALTER TABLE hydra_client ALTER COLUMN request_uris_json DROP DEFAULT;
-- hydra_client/post_logout_redirect_uris_json is meant to have a default
--split

ALTER TABLE hydra_client DROP COLUMN redirect_uris;
ALTER TABLE hydra_client DROP COLUMN grant_types;
ALTER TABLE hydra_client DROP COLUMN response_types;
ALTER TABLE hydra_client DROP COLUMN audience;
ALTER TABLE hydra_client DROP COLUMN allowed_cors_origins;
ALTER TABLE hydra_client DROP COLUMN contacts;
ALTER TABLE hydra_client DROP COLUMN request_uris;
ALTER TABLE hydra_client DROP COLUMN post_logout_redirect_uris;
--split

ALTER TABLE hydra_client RENAME COLUMN redirect_uris_json TO redirect_uris;
ALTER TABLE hydra_client RENAME COLUMN grant_types_json TO grant_types;
ALTER TABLE hydra_client RENAME COLUMN response_types_json TO response_types;
ALTER TABLE hydra_client RENAME COLUMN audience_json TO audience;
ALTER TABLE hydra_client RENAME COLUMN allowed_cors_origins_json TO allowed_cors_origins;
ALTER TABLE hydra_client RENAME COLUMN contacts_json TO contacts;
ALTER TABLE hydra_client RENAME COLUMN request_uris_json TO request_uris;
ALTER TABLE hydra_client RENAME COLUMN post_logout_redirect_uris_json TO post_logout_redirect_uris;
--split



ALTER TABLE hydra_oauth2_flow ADD COLUMN requested_scope_json jsonb NOT NULL DEFAULT '[]';
ALTER TABLE hydra_oauth2_flow ADD COLUMN requested_at_audience_json jsonb DEFAULT '[]';
ALTER TABLE hydra_oauth2_flow ADD COLUMN amr_json jsonb DEFAULT '[]';
ALTER TABLE hydra_oauth2_flow ADD COLUMN granted_scope_json jsonb;
ALTER TABLE hydra_oauth2_flow ADD COLUMN granted_at_audience_json jsonb DEFAULT '[]';
--split
UPDATE hydra_oauth2_flow SET requested_scope = '[]' WHERE requested_scope = '';
UPDATE hydra_oauth2_flow SET requested_at_audience = '[]' WHERE requested_at_audience = '';
UPDATE hydra_oauth2_flow SET amr = '[]' WHERE amr = '';
UPDATE hydra_oauth2_flow SET granted_scope = '[]' WHERE granted_scope = '';
UPDATE hydra_oauth2_flow SET granted_at_audience = '[]' WHERE granted_at_audience = '';
--split

UPDATE hydra_oauth2_flow SET requested_scope_json = cast('["' || REPLACE(requested_scope,'|','","') || '"]' as jsonb) WHERE requested_scope <> '[]';
UPDATE hydra_oauth2_flow SET requested_at_audience_json = cast('["' || REPLACE(requested_at_audience,'|','","') || '"]' as jsonb) WHERE requested_at_audience <> '[]';
UPDATE hydra_oauth2_flow SET amr_json = cast('["' || REPLACE(amr,'|','","') || '"]' as jsonb) WHERE amr <> '[]';
UPDATE hydra_oauth2_flow SET granted_scope_json = cast('["' || REPLACE(granted_scope,'|','","') || '"]' as jsonb) WHERE granted_scope <> '[]';
UPDATE hydra_oauth2_flow SET granted_at_audience_json = cast('["' || REPLACE(granted_at_audience,'|','","') || '"]' as jsonb) WHERE granted_at_audience <> '[]';
--split
ALTER TABLE hydra_oauth2_flow ALTER COLUMN requested_scope_json DROP DEFAULT;
--split

ALTER TABLE hydra_oauth2_flow DROP COLUMN requested_scope;
ALTER TABLE hydra_oauth2_flow DROP COLUMN requested_at_audience;
ALTER TABLE hydra_oauth2_flow DROP COLUMN amr;
ALTER TABLE hydra_oauth2_flow DROP COLUMN granted_scope;
ALTER TABLE hydra_oauth2_flow DROP COLUMN granted_at_audience;
--split

ALTER TABLE hydra_oauth2_flow RENAME COLUMN requested_scope_json TO requested_scope;
ALTER TABLE hydra_oauth2_flow RENAME COLUMN requested_at_audience_json TO requested_at_audience;
ALTER TABLE hydra_oauth2_flow RENAME COLUMN amr_json TO amr;
ALTER TABLE hydra_oauth2_flow RENAME COLUMN granted_scope_json TO granted_scope;
ALTER TABLE hydra_oauth2_flow RENAME COLUMN granted_at_audience_json TO granted_at_audience;

-- scripts/db-diff.sh can be used in code review to verify that the constraint hasn't changed; we need to recreate it due to the dropped and re-added columns
ALTER TABLE hydra_oauth2_flow ADD CHECK (
    state = 128 OR
    state = 129 OR
    state = 1 OR
    (state = 2 AND (
        login_remember IS NOT NULL AND
        login_remember_for IS NOT NULL AND
        login_error IS NOT NULL AND
        acr IS NOT NULL AND
        login_was_used IS NOT NULL AND
        context IS NOT NULL AND
        amr IS NOT NULL
    )) OR
    (state = 3 AND (
        login_remember IS NOT NULL AND
        login_remember_for IS NOT NULL AND
        login_error IS NOT NULL AND
        acr IS NOT NULL AND
        login_was_used IS NOT NULL AND
        context IS NOT NULL AND
        amr IS NOT NULL
    )) OR
    (state = 4 AND (
        login_remember IS NOT NULL AND
        login_remember_for IS NOT NULL AND
        login_error IS NOT NULL AND
        acr IS NOT NULL AND
        login_was_used IS NOT NULL AND
        context IS NOT NULL AND
        amr IS NOT NULL AND

        consent_challenge_id IS NOT NULL AND
        consent_verifier IS NOT NULL AND
        consent_skip IS NOT NULL AND
        consent_csrf IS NOT NULL
    )) OR
    (state = 5 AND (
        login_remember IS NOT NULL AND
        login_remember_for IS NOT NULL AND
        login_error IS NOT NULL AND
        acr IS NOT NULL AND
        login_was_used IS NOT NULL AND
        context IS NOT NULL AND
        amr IS NOT NULL AND

        consent_challenge_id IS NOT NULL AND
        consent_verifier IS NOT NULL AND
        consent_skip IS NOT NULL AND
        consent_csrf IS NOT NULL
    )) OR
    (state = 6 AND (
        login_remember IS NOT NULL AND
        login_remember_for IS NOT NULL AND
        login_error IS NOT NULL AND
        acr IS NOT NULL AND
        login_was_used IS NOT NULL AND
        context IS NOT NULL AND
        amr IS NOT NULL AND

        consent_challenge_id IS NOT NULL AND
        consent_verifier IS NOT NULL AND
        consent_skip IS NOT NULL AND
        consent_csrf IS NOT NULL AND

        granted_scope IS NOT NULL AND
        consent_remember IS NOT NULL AND
        consent_remember_for IS NOT NULL AND
        consent_error IS NOT NULL AND
        session_access_token IS NOT NULL AND
        session_id_token IS NOT NULL AND
        consent_was_used IS NOT NULL
    ))
)
