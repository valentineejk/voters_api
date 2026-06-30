CREATE TABLE users (
    id           TEXT        PRIMARY KEY,      -- USR-00001
    email        TEXT        NOT NULL UNIQUE,
    password_hash TEXT       NOT NULL,         -- bcrypt hash, never plain text
    role         TEXT        NOT NULL DEFAULT 'user',  -- user | admin
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id          TEXT        PRIMARY KEY,       -- random UUID
    user_id     TEXT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT        NOT NULL UNIQUE,   -- hash of the token, not the token itself
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked     BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);