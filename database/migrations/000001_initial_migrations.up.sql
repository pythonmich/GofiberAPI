SET TIMEZONE = 'Africa/Nairobi';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users(
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    password_changed_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE UNIQUE INDEX user_email ON users(email);
