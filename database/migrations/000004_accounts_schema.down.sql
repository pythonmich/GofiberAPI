ALTER TABLE accounts DROP CONSTRAINT IF EXISTS "accounts_user_id_fkey";
ALTER TABLE accounts DROP CONSTRAINT IF EXISTS "owner_accounts_uiq";

DROP TABLE accounts CASCADE ;