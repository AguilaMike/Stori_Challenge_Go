CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES accounts(id),
    amount DECIMAL(15, 2) NOT NULL,
    type TEXT NOT NULL,
    input_file_id TEXT NOT NULL,
    input_date TIMESTAMP NOT NULL,
    created_at BIGINT NOT NULL
);
