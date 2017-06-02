package main

import (
	"api"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func CreateGameHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		logger := log.With(Ctx(r).Logger, "handler", "create game")

		res, err := db.Exec("INSERT INTO game (created_at) VALUES (?)", time.Now())
		if err != nil {
			logger.Log("lvl", "error", "msg", "creating new game", "err", err.Error())
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			logger.Log("lvl", "error", "msg", "retrieving game id", "err", err.Error())
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		logger.Log("lvl", "error", "msg", "created new game")
		api.Write(w, map[string]interface{}{
			"game": id,
		})
	}
}
