ALTER TABLE sessions DROP CONSTRAINT IF EXISTS "sessions_user_id_fkey";
ALTER TABLE sessions DROP CONSTRAINT IF EXISTS "current_sessions_uiq";

DROP TABLE IF EXISTS sessions CASCADE ;