CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  fullname varchar(30) NOT NULL,
  email citext UNIQUE NOT NULL,
  phone VARCHAR(20),
  slug varchar(50) NOT NULL,
  dob varchar(20) NOT NULL,
  gender varchar(10) NOT NULL,
  address varchar(30) NOT NULL,
  occupation varchar(30),
  user_type varchar(25),
  id_type varchar(30),
  id_number varchar(30),
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  is_active BOOLEAN NOT NULL DEFAULT FALSE
);