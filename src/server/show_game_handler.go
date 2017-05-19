package main

import (
	"api"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func ShowGameHandler(logger log.Logger, db *sqlx.DB) httprouter.Handle {
	logger = log.With(logger, "handler", "show game")

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			api.WriteError(w, http.StatusBadRequest, err)
			return
		}

		var game Game
		err = db.Get(&game, "SELECT id, created_at, started_at, finished_at FROM game WHERE id = ?", id)
		if err == sql.ErrNoRows {
			api.WriteError(w, http.StatusNotFound, err)
			return
		}
		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		api.Write(w, game)
	}
}
