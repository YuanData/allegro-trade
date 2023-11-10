CREATE TABLE "members" (
  "membername" varchar PRIMARY KEY,
  "password_hash" varchar NOT NULL,
  "name_entire" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_time" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),  
  "created_time" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "traders" ADD FOREIGN KEY ("holder") REFERENCES "members" ("membername");

ALTER TABLE "traders" ADD CONSTRAINT "holder_symbol_key" UNIQUE ("holder", "symbol");
