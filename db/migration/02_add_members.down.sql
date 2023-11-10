ALTER TABLE IF EXISTS "traders" DROP CONSTRAINT IF EXISTS "holder_symbol_key";

ALTER TABLE IF EXISTS "traders" DROP CONSTRAINT IF EXISTS "traders_holder_fkey";

DROP TABLE IF EXISTS "members";
