package main

import (
	"api"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func CreateGameHandler(logger log.Logger, db *sqlx.DB) httprouter.Handle {
	logger = log.With(logger, "handler", "create game")

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		res, err := db.Exec("INSERT INTO game (created_at) VALUES (?)", time.Now())
		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		api.Write(w, map[string]interface{}{
			"game": id,
		})
	}
}
