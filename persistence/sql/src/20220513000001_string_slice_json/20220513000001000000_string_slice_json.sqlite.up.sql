UPDATE hydra_client SET redirect_uris = '[]' WHERE redirect_uris = '';
UPDATE hydra_client SET grant_types = '[]' WHERE grant_types = '';
UPDATE hydra_client SET response_types = '[]' WHERE response_types = '';
UPDATE hydra_client SET audience = '[]' WHERE audience = '';
UPDATE hydra_client SET allowed_cors_origins = '[]' WHERE allowed_cors_origins = '';
UPDATE hydra_client SET contacts = '[]' WHERE contacts = '';
UPDATE hydra_client SET request_uris = '[]' WHERE request_uris = '';
UPDATE hydra_client SET post_logout_redirect_uris = '[]' WHERE post_logout_redirect_uris = '';
--split

UPDATE hydra_client SET redirect_uris = '["' || REPLACE(redirect_uris,'|','","') || '"]' WHERE redirect_uris <> '[]';
UPDATE hydra_client SET grant_types = '["' || REPLACE(grant_types,'|','","') || '"]' WHERE grant_types <> '[]';
UPDATE hydra_client SET response_types = '["' || REPLACE(response_types,'|','","') || '"]' WHERE response_types <> '[]';
UPDATE hydra_client SET audience = '["' || REPLACE(audience,'|','","') || '"]' WHERE audience <> '[]';
UPDATE hydra_client SET allowed_cors_origins = '["' || REPLACE(allowed_cors_origins,'|','","') || '"]' WHERE allowed_cors_origins <> '[]';
UPDATE hydra_client SET contacts = '["' || REPLACE(contacts,'|','","') || '"]' WHERE contacts <> '[]';
UPDATE hydra_client SET request_uris = '["' || REPLACE(request_uris,'|','","') || '"]' WHERE request_uris <> '[]';
UPDATE hydra_client SET post_logout_redirect_uris = '["' || REPLACE(post_logout_redirect_uris,'|','","') || '"]' WHERE post_logout_redirect_uris <> '[]';
--split

UPDATE hydra_oauth2_flow SET requested_scope = '[]' WHERE requested_scope = '';
UPDATE hydra_oauth2_flow SET requested_at_audience = '[]' WHERE requested_at_audience = '';
UPDATE hydra_oauth2_flow SET amr = '[]' WHERE amr = '';
UPDATE hydra_oauth2_flow SET granted_scope = '[]' WHERE granted_scope = '';
UPDATE hydra_oauth2_flow SET granted_at_audience = '[]' WHERE granted_at_audience = '';
--split

UPDATE hydra_oauth2_flow SET requested_scope = '["' || REPLACE(requested_scope,'|','","') || '"]' WHERE requested_scope <> '[]';
UPDATE hydra_oauth2_flow SET requested_at_audience = '["' || REPLACE(requested_at_audience,'|','","') || '"]' WHERE requested_at_audience <> '[]';
UPDATE hydra_oauth2_flow SET amr = '["' || REPLACE(amr,'|','","') || '"]' WHERE amr <> '[]';
UPDATE hydra_oauth2_flow SET granted_scope = '["' || REPLACE(granted_scope,'|','","') || '"]' WHERE granted_scope <> '[]';
UPDATE hydra_oauth2_flow SET granted_at_audience = '["' || REPLACE(granted_at_audience,'|','","') || '"]' WHERE granted_at_audience <> '[]';