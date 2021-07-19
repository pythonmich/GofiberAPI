-- create type roles to avoid incorrect input
-- we only need an admin role only
-- Role 'member' is if user exists in database
-- futures roles if needed will be added to this enum


CREATE TYPE user_role AS ENUM (
        'admin'
);

CREATE TABLE IF NOT EXISTS user_roles(
    user_id UUID NOT NULL REFERENCES users,
    role user_role NOT NULL ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    PRIMARY KEY (user_id, role)
);

CREATE UNIQUE INDEX user_roles_uiq ON user_roles(user_id)