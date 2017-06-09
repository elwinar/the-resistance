package main

import "time"

type Player struct {
	ID       int       `db:"id"        json:"id"`
	UserID   int       `db:"user_id"   json:"user_id"`
	Name     string    `db:"name"      json:"name"`
	JoinedAt time.Time `db:"joined_at" json:"joined_at"`
}
