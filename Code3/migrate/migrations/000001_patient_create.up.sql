CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS patients (
  id bigserial PRIMARY KEY,
  fullname text NOT NULL,
  phone varchar(100),
  email citext UNIQUE NOT NULL,
  dob varchar(20),
  slug varchar(50) NOT NULL,
  gender varchar(10) NOT NULL,
  address varchar(30) NOT NULL,
  occupation varchar(30),
  emergency_name varchar(30),
  emergency_phone varchar(20),
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  is_active BOOLEAN NOT NULL DEFAULT FALSE
);