package main

import (
	"api"
	"net/http"

	"github.com/elwinar/token"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func AuthenticateHandler(logger log.Logger, db *sqlx.DB, secret []byte) httprouter.Handle {
	logger = log.With(logger, "handler", "authenticate")

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		claims, err := token.ParseHS256(r.Header.Get("token"), secret)
		if err != nil {
			api.Write(w, map[string]interface{}{
				"authenticated": false,
			})
			return
		}

		api.Write(w, map[string]interface{}{
			"token":         r.Header.Get("token"),
			"authenticated": true,
		})
	}
}
