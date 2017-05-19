package main

import (
	"api"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func ListGameHandler(logger log.Logger, db *sqlx.DB) httprouter.Handle {
	logger = log.With(logger, "handler", "show game")

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var games []Game
		err := db.Select(&games, "SELECT id, created_at, started_at, finished_at FROM game")
		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		api.Write(w, games)
	}
}
