CREATE TABLE IF NOT EXISTS cards
	(
		card_id SERIAL PRIMARY KEY,
		user_id BIGINT REFERENCES users (user_id) NOT NULL,
		card_number VARCHAR(16) NOT NULL,
		card_holder VARCHAR(255) NOT NULL,
		card_expiration_date DATE NOT NULL,
		card_cvv VARCHAR(3) NOT NULL,
		card_bank VARCHAR(255) NOT NULL,
		metadata JSONB NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NULL,
		deleted_at TIMESTAMP NULL
	);
CREATE INDEX idx_card_user_id ON cards (user_id);
