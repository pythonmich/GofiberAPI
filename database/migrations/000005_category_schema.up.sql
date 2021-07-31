CREATE TABLE IF NOT EXISTS categories(
     category_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
     parent_id VARCHAR NOT NULL DEFAULT uuid_nil(),
     user_id UUID NOT NULL REFERENCES users,
     name VARCHAR NOT NULL,
     created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
     deleted_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

ALTER TABLE categories ADD CONSTRAINT "user_category_uiq" UNIQUE (parent_id, user_id, name)
