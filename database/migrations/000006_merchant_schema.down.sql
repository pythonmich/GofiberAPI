ALTER TABLE merchant DROP CONSTRAINT IF EXISTS "merchant_user_id_fkey";
ALTER TABLE merchant DROP CONSTRAINT IF EXISTS "user_merchant_uiq";

DROP TABLE IF EXISTS merchant;