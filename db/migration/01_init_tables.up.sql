CREATE TABLE "traders" (
  "id" bigserial PRIMARY KEY,
  "holder" varchar NOT NULL,
  "rest" bigint NOT NULL,
  "symbol" varchar NOT NULL,
  "created_time" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "records" (
  "id" bigserial PRIMARY KEY,
  "trader_id" bigint NOT NULL,
  "collection" bigint NOT NULL,
  "created_time" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "details" (
  "id" bigserial PRIMARY KEY,
  "from_trader_id" bigint NOT NULL,
  "to_trader_id" bigint NOT NULL,
  "collection" bigint NOT NULL,
  "created_time" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "records" ADD FOREIGN KEY ("trader_id") REFERENCES "traders" ("id");

ALTER TABLE "details" ADD FOREIGN KEY ("from_trader_id") REFERENCES "traders" ("id");

ALTER TABLE "details" ADD FOREIGN KEY ("to_trader_id") REFERENCES "traders" ("id");

CREATE INDEX ON "traders" ("holder");

CREATE INDEX ON "records" ("trader_id");

CREATE INDEX ON "details" ("from_trader_id");

CREATE INDEX ON "details" ("to_trader_id");

CREATE INDEX ON "details" ("from_trader_id", "to_trader_id");

COMMENT ON COLUMN "records"."collection" IS 'can be negative or positive';

COMMENT ON COLUMN "details"."collection" IS 'must be positive';
