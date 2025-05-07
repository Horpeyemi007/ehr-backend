-- Roles table
CREATE TABLE IF NOT EXISTS roles (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description TEXT DEFAULT NULL,
  type VARCHAR(20) NOT NULL
);

-- Permissions table
CREATE TABLE IF NOT EXISTS permissions (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description TEXT DEFAULT NULL,
  type VARCHAR(20) NOT NULL
);

-- Create User-Roles relation table
CREATE TABLE user_roles (
	user_id BIGSERIAL,
	role_id BIGSERIAL,
  description TEXT DEFAULT NULL,
	PRIMARY KEY (user_id, role_id),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- Role-Permissions mapping
CREATE TABLE IF NOT EXISTS role_permissions (
  role_id BIGSERIAL NOT NULL,
  permission_id BIGSERIAL NOT NULL,
  description TEXT DEFAULT NULL,
  PRIMARY KEY (role_id,permission_id),
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- User-Permissions mapping (direct user-level permissions)
CREATE TABLE IF NOT EXISTS user_permissions (
  user_id BIGSERIAL NOT NULL,
  permission_id BIGSERIAL NOT NULL,
  description TEXT DEFAULT NULL,
  PRIMARY KEY (user_id, permission_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);