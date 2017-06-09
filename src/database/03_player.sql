CREATE TABLE player (
	id INTEGER PRIMARY KEY,
	game_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	name VARCHAR(50) NOT NULL,
	joined_at DATETIME NOT NULL,
	FOREIGN KEY (user_id) REFERENCES user(id),
	FOREIGN KEY (game_id) REFERENCES game(id),
	UNIQUE (game_id, user_id)
);
