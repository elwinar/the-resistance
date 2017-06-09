CREATE TABLE player (
	id INTEGER PRIMARY KEY,
	game_id INTEGER NOT NULL,
	player_id INTEGER NOT NULL,
	name VARCHAR(50) NOT NULL,
	joined_at DATETIME NOT NULL,
	FOREIGN KEY (player_id) REFERENCES user(id),
	FOREIGN KEY (game_id) REFERENCES game(id),
	UNIQUE (game_id, player_id)
);
