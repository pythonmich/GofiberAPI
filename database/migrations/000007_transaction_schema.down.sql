ALTER TABLE transactions DROP CONSTRAINT IF EXISTS "transactions_user_id_fkey";
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS "transactions_account_id_fkey";
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS "transactions_category_id_fkey";
ALTER TABLE accounts ALTER COLUMN account_type TYPE varchar USING (account_type::varchar::accounts_type);


DROP TABLE IF EXISTS transactions;
DROP TYPE IF EXISTS accounts_type;
DROP TYPE IF EXISTS transactions_type;