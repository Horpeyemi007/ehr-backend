CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS admins (
  id bigserial PRIMARY KEY,
  fullname VARCHAR(40) NOT NULL,
  email citext UNIQUE NOT NULL,
  password bytea NOT NULL,
  slug varchar(30) NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  is_active BOOLEAN NOT NULL DEFAULT FALSE
);
