-- migrate:up
CREATE TABLE point_of_interest (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    updated_at TIMESTAMPTZ,
    updated_by TEXT,
    deleted_at TIMESTAMPTZ,
    deleted_by TEXT
);

-- migrate:down
DROP TABLE IF EXISTS point_of_interest;