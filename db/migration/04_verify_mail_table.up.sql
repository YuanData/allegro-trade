CREATE TABLE "verify_emails" (
  "id" bigserial PRIMARY KEY,
  "membername" varchar NOT NULL,
  "email" varchar NOT NULL,
  "secret_code" varchar NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "created_time" timestamptz NOT NULL DEFAULT (now()),
  "expired_time" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("membername") REFERENCES "members" ("membername");

ALTER TABLE "members" ADD COLUMN "is_email_verified" bool NOT NULL DEFAULT false;
