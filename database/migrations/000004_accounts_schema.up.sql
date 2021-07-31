CREATE TABLE accounts (
    account_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users,
    account_name VARCHAR NOT NULL,
    account_type VARCHAR NOT NULL,
    balance BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    deleted_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

ALTER TABLE accounts ADD CONSTRAINT "owner_accounts_uiq" UNIQUE (user_id, currency, account_name, account_type)