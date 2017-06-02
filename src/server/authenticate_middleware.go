package main

import (
	"api"
	"net/http"

	"github.com/elwinar/token"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/negroni"
)

func AuthenticateMiddleware(logger log.Logger, db *sqlx.DB, secret []byte, unprotectedPaths []string) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		logger := log.With(Ctx(r).Logger, "handler", "authenticate")

		for _, path := range unprotectedPaths {
			if r.URL.Path == path {
				next(w, r)
				return
			}
		}

		_, err := token.ParseHS256(r.Header.Get("token"), secret)
		if err != nil {
			logger.Log("lvl", "error", "msg", "authentication failed", "err", err.Error())
			api.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		next(w, r)
	}
}
