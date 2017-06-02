package main

type User struct {
	ID       uint64 `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}
