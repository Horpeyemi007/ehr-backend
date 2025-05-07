CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS entity (
  id bigserial PRIMARY KEY,
  fullname VARCHAR(40) NOT NULL,
  email citext UNIQUE NOT NULL,
  password bytea NOT NULL,
  user_type VARCHAR(20) NOT NULL DEFAULT 'root',
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  is_active BOOLEAN NOT NULL DEFAULT FALSE
);


CREATE TABLE IF NOT EXISTS entity_information (
  slug 	VARCHAR(30) NOT NULL,
  entity_id INT NOT NULL,
  FOREIGN KEY (entity_id) REFERENCES entity(id),
	PRIMARY KEY (slug)
);