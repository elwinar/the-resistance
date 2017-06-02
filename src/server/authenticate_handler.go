package main

import (
	"api"
	"net/http"

	"github.com/elwinar/token"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func AuthenticateHandler(db *sqlx.DB, secret []byte) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ctx := Ctx(r)
		logger := log.With(ctx.Logger, "handler", "authenticate")

		t := r.Header.Get("token")
		_, err := token.ParseHS256(t, secret)
		if err != nil {
			logger.Log("lvl", "error", "msg", "invalid token", "err", err.Error(), "token", t)
			api.Write(w, map[string]interface{}{
				"authenticated": false,
			})
			return
		}

		logger.Log("lvl", "info", "msg", "valid token", "token", t)
		api.Write(w, map[string]interface{}{
			"authenticated": true,
		})
	}
}
