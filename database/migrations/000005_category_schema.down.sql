ALTER TABLE categories DROP CONSTRAINT IF EXISTS "user_category_uiq";
ALTER TABLE categories DROP CONSTRAINT IF EXISTS "categories_parent_id_fkey";

DROP TABLE IF EXISTS categories CASCADE ;