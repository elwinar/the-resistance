package main

import "time"

type Player struct {
	ID       int       `db:"id"        json:"id"`
	Name     string    `db:"name"      json:"name"`
	JoinedAt time.Time `db:"joined_at" json:"joined_at"`
}
