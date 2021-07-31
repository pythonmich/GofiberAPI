CREATE TABLE IF NOT EXISTS merchant (
    merchant_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users,
    name VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    deleted_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

ALTER TABLE merchant ADD CONSTRAINT "user_merchant_uiq" UNIQUE (user_id, name)


