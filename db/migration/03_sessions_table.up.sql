CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "membername" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expired_time" timestamptz NOT NULL,
  "created_time" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("membername") REFERENCES "members" ("membername");
