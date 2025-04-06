CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS patient (
  id bigserial PRIMARY KEY,
  fullname text NOT NULL,
  dob varchar(20),
  sex varchar(15) NOT NULL,
  phone varchar(100),
  email citext UNIQUE NOT NULL,
  street text,
  city varchar(30),
  state varchar(20),
  zipcode varchar(20),
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);