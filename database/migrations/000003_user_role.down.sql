ALTER TABLE IF EXISTS users_roles DROP CONSTRAINT IF EXISTS "user_roles_uiq";
ALTER TABLE IF EXISTS users_roles DROP CONSTRAINT IF EXISTS "user_roles_user_id_fkey";

DROP TABLE IF EXISTS user_roles;
DROP TYPE IF EXISTS user_role;
