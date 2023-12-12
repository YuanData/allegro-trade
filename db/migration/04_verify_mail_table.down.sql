DROP TABLE IF EXISTS "verify_emails" CASCADE;

ALTER TABLE "members" DROP COLUMN "is_email_verified";
