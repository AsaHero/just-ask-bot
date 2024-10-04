CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users"(
  "guid" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "external_id" bigint NOT NULL,
  "username" varchar,
  "first_name" varchar,
  "last_name" varchar,
  "is_blocked" bool DEFAULT FALSE,
  "created_at" timestamp with time zone DEFAULT now(),
  "updated_at" timestamp with time zone
);

CREATE INDEX idx_users_external_id ON "users"("external_id");

CREATE INDEX idx_users_username ON "users"("username");

