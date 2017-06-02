CREATE TABLE game_joined (
	game_id INTEGER NOT NULL,
	player_id INTEGER NOT NULL,
	joined_at DATETIME NOT NULL,
	FOREIGN KEY (player_id) REFERENCES user(id),
	FOREIGN KEY (game_id) REFERENCES game(id),
	PRIMARY KEY (game_id, player_id)
);
