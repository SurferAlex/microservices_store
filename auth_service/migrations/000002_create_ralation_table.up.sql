CREATE TABLE IF NOT EXISTS role_permissions (
	      role_id INT NOT NULL,
		  permission_id INT NOT NULL,
		  PRIMARY KEY (role_id, permission_id),
		  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
		  FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
	);

CREATE TABLE IF NOT EXISTS user_roles (
		user_id INT NOT NULL,
		role_id INT NOT NULL,
		PRIMARY KEY (user_id, role_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
	);

CREATE TABLE IF NOT EXISTS refresh_tokens (
		token_hash VARCHAR(64) PRIMARY KEY,
		user_id INT NOT NULL,
		expires_at TIMESTAMPTZ NOT NULL,
		revoked_at TIMESTAMPTZ,
		user_agent TEXT,
		ip TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);    