CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    nickname TEXT NOT NULL,
    email TEXT NOT NULL,
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE
);
