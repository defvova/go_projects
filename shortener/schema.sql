CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email text NOT NULL,
    password_hash text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_email ON users(email);

-- DROP TABLE users;

CREATE TABLE IF NOT EXISTS sessions (
    id BIGSERIAL PRIMARY KEY,
    expires_at TIMESTAMPTZ NOT NULL,
    user_id BIGINT NOT NULL,
    token_hash BYTEA NOT NULL,
    CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_session_token_hash ON sessions(token_hash);

-- DROP TABLE sessions;

CREATE TABLE IF NOT EXISTS redirects (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    url text NOT NULL,
    token text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_redirects_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_redirect_user_id ON redirects(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_redirect_token ON redirects(token);

-- DROP TABLE redirects;
