package main

import "time"

type Game struct {
	ID         int        `db:"id"          json:"id"`
	Name       string     `db:"name"        json:"name"`
	Players    int        `db:"players"     json:"players"`
	CreatedAt  time.Time  `db:"created_at"  json:"created_at"`
	StartedAt  *time.Time `db:"started_at"  json:"started_at"`
	FinishedAt *time.Time `db:"finished_at" json:"finished_at"`
}
