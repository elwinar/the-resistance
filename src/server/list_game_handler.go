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
		ctx := Ctx(r)
		logger := log.With(ctx.Logger, "handler", "list games")

		var games []Game
		err := db.Select(&games, `
			SELECT g.id, g.name, g.players, COALESCE(p.joined, 0) AS joined, g.created_at, g.started_at, g.finished_at 
			FROM game AS g 
			LEFT OUTER JOIN (
				SELECT game_id, COUNT(*) AS joined
				FROM player 
				GROUP BY game_id
			) AS p ON g.id = p.game_id
		`)
		if err != nil {
			logger.Log("lvl", "error", "msg", "retrieving games list", "err", err.Error())
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		api.Write(w, games)
	}
}
