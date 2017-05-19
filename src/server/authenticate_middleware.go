package main

import (
	"api"
	"net/http"

	"github.com/elwinar/token"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/negroni"
)

func AuthenticateMiddleware(logger log.Logger, db *sqlx.DB, secret []byte) negroni.HandlerFunc {
	logger = log.With(logger, "handler", "authenticate")

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		//
		if r.URL.Path == "/login" {
			next(w, r)
			return
		}

		_, err := token.ParseHS256(r.Header.Get("token"), secret)
		if err != nil {
			api.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		next(w, r)
	}
}
