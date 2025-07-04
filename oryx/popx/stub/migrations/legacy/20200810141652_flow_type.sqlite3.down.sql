CREATE TABLE "_selfservice_login_requests_tmp" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"active_method" TEXT NOT NULL,
"csrf_token" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"forced" bool NOT NULL DEFAULT 'false',
"messages" TEXT
);
INSERT INTO "_selfservice_login_requests_tmp" (id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at, forced, messages) SELECT id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at, forced, messages FROM "selfservice_login_requests";

DROP TABLE "selfservice_login_requests";
ALTER TABLE "_selfservice_login_requests_tmp" RENAME TO "selfservice_login_requests";
CREATE TABLE "_selfservice_registration_requests_tmp" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"active_method" TEXT NOT NULL,
"csrf_token" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"messages" TEXT
);
INSERT INTO "_selfservice_registration_requests_tmp" (id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at, messages) SELECT id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at, messages FROM "selfservice_registration_requests";

DROP TABLE "selfservice_registration_requests";
ALTER TABLE "_selfservice_registration_requests_tmp" RENAME TO "selfservice_registration_requests";
CREATE TABLE "_selfservice_settings_requests_tmp" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"identity_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"active_method" TEXT,
"messages" TEXT,
"state" TEXT NOT NULL DEFAULT 'show_form',
FOREIGN KEY (identity_id) REFERENCES identities (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
INSERT INTO "_selfservice_settings_requests_tmp" (id, request_url, issued_at, expires_at, identity_id, created_at, updated_at, active_method, messages, state) SELECT id, request_url, issued_at, expires_at, identity_id, created_at, updated_at, active_method, messages, state FROM "selfservice_settings_requests";

DROP TABLE "selfservice_settings_requests";
ALTER TABLE "_selfservice_settings_requests_tmp" RENAME TO "selfservice_settings_requests";
CREATE TABLE "_selfservice_recovery_requests_tmp" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"messages" TEXT,
"active_method" TEXT,
"csrf_token" TEXT NOT NULL,
"state" TEXT NOT NULL,
"recovered_identity_id" char(36),
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (recovered_identity_id) REFERENCES identities (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
INSERT INTO "_selfservice_recovery_requests_tmp" (id, request_url, issued_at, expires_at, messages, active_method, csrf_token, state, recovered_identity_id, created_at, updated_at) SELECT id, request_url, issued_at, expires_at, messages, active_method, csrf_token, state, recovered_identity_id, created_at, updated_at FROM "selfservice_recovery_requests";

DROP TABLE "selfservice_recovery_requests";
ALTER TABLE "_selfservice_recovery_requests_tmp" RENAME TO "selfservice_recovery_requests";
CREATE TABLE "_selfservice_verification_requests_tmp" (
"id" TEXT PRIMARY KEY,
"request_url" TEXT NOT NULL,
"issued_at" DATETIME NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
"expires_at" DATETIME NOT NULL,
"csrf_token" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"messages" TEXT,
"via" TEXT NOT NULL DEFAULT 'email',
"success" bool NOT NULL DEFAULT 'FALSE'
);
INSERT INTO "_selfservice_verification_requests_tmp" (id, request_url, issued_at, expires_at, csrf_token, created_at, updated_at, messages, via, success) SELECT id, request_url, issued_at, expires_at, csrf_token, created_at, updated_at, messages, via, success FROM "selfservice_verification_requests";

DROP TABLE "selfservice_verification_requests";
ALTER TABLE "_selfservice_verification_requests_tmp" RENAME TO "selfservice_verification_requests";