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

func ShowGameHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := Ctx(r)
		logger := log.With(ctx.Logger, "handler", "show game")

		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			logger.Log("lvl", "error", "msg", "parsing game id", "err", err.Error())
			api.WriteError(w, http.StatusBadRequest, err)
			return
		}
		logger = log.With(logger, "game", id)

		var game Game
		err = db.Get(&game, "SELECT g.id, g.name, g.players, COUNT(p.id) AS joined, g.created_at, g.started_at, g.finished_at FROM game AS g LEFT OUTER JOIN player AS p ON g.id = p.game_id WHERE g.id = ?", id)
		if err == sql.ErrNoRows {
			logger.Log("lvl", "error", "msg", "game not found", "err", err.Error())
			api.WriteError(w, http.StatusNotFound, err)
			return
		}
		if err != nil {
			logger.Log("lvl", "error", "msg", "retrieving game", "err", err.Error())
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		logger.Log("lvl", "info", "msg", "retrieved game")
		api.Write(w, game)
	}
}
