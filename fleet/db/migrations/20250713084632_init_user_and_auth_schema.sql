-- migrate:up
CREATE TABLE "user" (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    email text NOT NULL,
    email_hash text NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMPTZ NULL,
    updated_by VARCHAR(255) NULL,
    deleted_at TIMESTAMPTZ NULL,
    deleted_by VARCHAR(255) NULL
);

-- Create trigger function for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to user table
CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- AUTH TABLE
CREATE TABLE auth (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    password text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMPTZ NULL,
    updated_by VARCHAR(255) NULL,
    deleted_at TIMESTAMPTZ NULL,
    deleted_by VARCHAR(255) NULL,
    CONSTRAINT fk_auth_user_id FOREIGN KEY (user_id) REFERENCES "user" (id)
);

-- Apply trigger to auth table
CREATE TRIGGER update_auth_updated_at
    BEFORE UPDATE ON auth
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- SESSION TABLE
CREATE TABLE session (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMPTZ NULL,
    updated_by VARCHAR(255) NULL,
    deleted_at TIMESTAMPTZ NULL,
    deleted_by VARCHAR(255) NULL,
    CONSTRAINT fk_session_user_id FOREIGN KEY (user_id) REFERENCES "user" (id)
);

-- Apply trigger to session table
CREATE TRIGGER update_session_updated_at
    BEFORE UPDATE ON session
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- INDEXES
CREATE INDEX idx_user_email_hash_deleted_at ON "user" (email_hash, deleted_at);
CREATE INDEX idx_auth_user_id_deleted_at ON auth (user_id, deleted_at);
CREATE INDEX idx_session_user_id ON session (user_id);
CREATE INDEX idx_session_expires_at ON session (expires_at);
CREATE INDEX idx_session_id ON session (id);

-- migrate:down
-- Drop triggers first
DROP TRIGGER IF EXISTS update_session_updated_at ON session;
DROP TRIGGER IF EXISTS update_auth_updated_at ON auth;
DROP TRIGGER IF EXISTS update_user_updated_at ON "user";

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (to handle dependencies)
DROP TABLE IF EXISTS session CASCADE;
DROP TABLE IF EXISTS auth CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
