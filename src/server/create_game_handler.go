package main

import (
	"api"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func CreateGameHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ctx := Ctx(r)
		logger := log.With(ctx.Logger, "handler", "create game")

		var game Game
		err := json.NewDecoder(r.Body).Decode(&game)
		if err != nil {
			logger.Log("lvl", "error", "msg", "parsing payload", "err", err.Error())
			api.WriteError(w, http.StatusBadRequest, err)
			return
		}

		if game.Players < 5 || game.Players > 10 {
			err = errors.New("invalid number of players (5 =< x <= 10)")
			logger.Log("lvl", "error", "msg", "invalid number of players", "err", err.Error())
			api.WriteError(w, http.StatusBadRequest, err)
			return
		}

		res, err := db.Exec("INSERT INTO game (name, players, created_at) VALUES (?, ?, ?)", game.Name, game.Players, time.Now())
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

		logger.Log("lvl", "info", "msg", "created new game")
		api.Write(w, map[string]interface{}{
			"game": id,
		})
	}
}
