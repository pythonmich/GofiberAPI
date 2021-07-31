-- other types shall be added this are examples

CREATE TYPE accounts_type AS ENUM (
    'cash',
    'credit'
);

CREATE TYPE transactions_type AS ENUM (
    'income',
    'expense'
);


CREATE TABLE IF NOT EXISTS transactions(
    transaction_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users,
    account_id UUID NOT NULL REFERENCES accounts,
    category_id UUID NOT NULL REFERENCES categories,
    name VARCHAR NOT NULL,
    transaction_type transactions_type NOT NULL,
    amount INTEGER NOT NULL,
    notes VARCHAR NOT NULL DEFAULT '',
    date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    deleted_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

ALTER TABLE accounts ALTER COLUMN account_type TYPE accounts_type USING (trim(account_type):: accounts_type);
ALTER TABLE categories ALTER parent_id DROP DEFAULT;
ALTER TABLE categories ALTER COLUMN parent_id SET DEFAULT '';