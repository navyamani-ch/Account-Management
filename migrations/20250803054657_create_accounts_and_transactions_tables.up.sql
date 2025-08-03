CREATE TABLE accounts (
                          id SERIAL PRIMARY KEY,
                          account_id INTEGER NOT NULL UNIQUE,
                          balance NUMERIC NOT NULL DEFAULT 0,
                          created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE transactions (
                              id SERIAL PRIMARY KEY,
                              source_account_id INTEGER NOT NULL REFERENCES accounts(account_id) ON DELETE CASCADE,
                              destination_account_id INTEGER NOT NULL REFERENCES accounts(account_id) ON DELETE CASCADE,
                              amount NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
