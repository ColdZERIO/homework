CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS user_token (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    refresh_token TEXT NOT NULL UNIQUE DEFAULT '',
    expired_at TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL
);
