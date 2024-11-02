CREATE TABLE IF NOT EXISTS users
	  (
      user_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
      username VARCHAR(255) NOT NULL UNIQUE,
      password VARCHAR(255) NOT NULL,
      created_at TIMESTAMP NOT NULL,
      updated_at TIMESTAMP NULL,
      deleted_at TIMESTAMP NULL
    );
    CREATE UNIQUE INDEX idx_username_unique ON users (username);

	CREATE TABLE IF NOT EXISTS entities
		(
		  entity_id SERIAL PRIMARY KEY,
		  user_id BIGINT REFERENCES users (user_id) NOT NULL,
		  metadata JSONB NOT NULL,
		  created_at TIMESTAMP NOT NULL,
		  updated_at TIMESTAMP NULL,
		  deleted_at TIMESTAMP NULL
		);
CREATE INDEX idx_entity_user_id ON entities (user_id);
