package main

import (
	"api"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func JoinGameHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := Ctx(r)
		logger := log.With(ctx.Logger, "handler", "join game")

		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			logger.Log("lvl", "error", "msg", "parsing game id", "err", err.Error())
			api.WriteError(w, http.StatusBadRequest, err)
			return
		}
		logger = log.With(logger, "game", id)

		var game Game
		err = db.Get(&game, "SELECT id, name, players, created_at, started_at, finished_at FROM game WHERE id = ?", id)
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

		var count int
		err = db.Get(&count, "SELECT count(*) FROM game_joined WHERE game_id = ?", game.ID)
		if err != nil {
			logger.Log("lvl", "error", "msg", "counting players who joined game", "err", err.Error())
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if count >= game.Players {
			err = errors.New("game already full")
			logger.Log("lvl", "error", "msg", "can't add player to game", "err", err.Error())
			api.WriteError(w, http.StatusBadRequest, err)
			return
		}

		_, err = db.Exec("INSERT INTO game_joined (game_id, player_id, joined_at) VALUES (?, ?, ?)", game.ID, ctx.UserID, time.Now())
		if err != nil {
			logger.Log("lvl", "error", "msg", "joining game", "err", err.Error())
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		logger.Log("lvl", "info", "msg", "joined game", "game", game.ID)
		api.Write(w, "ok")
	}
}