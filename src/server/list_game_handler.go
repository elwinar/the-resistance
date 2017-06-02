package main

import (
	"api"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func ListGameHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		logger := log.With(Ctx(r).Logger, "handler", "list games")

		var games []Game
		err := db.Select(&games, "SELECT id, created_at, started_at, finished_at FROM game")
		if err != nil {
			logger.Log("lvl", "error", "msg", "retrieving games list", "err", err.Error())
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		api.Write(w, games)
	}
}
