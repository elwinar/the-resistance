package main

import (
	"api"
	"errors"
	"net/http"
	"time"

	"github.com/elwinar/token"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func LoginHandler(logger log.Logger, db *sqlx.DB, secret []byte, tokenTTL time.Duration) httprouter.Handle {
	logger = log.With(logger, "handler", "login")

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var req LoginRequest
		err := api.Read(r, &req)
		if err != nil {
			logger.Log("lvl", "error", "msg", "reading request", "err", err)
			api.WriteError(w, 400, err)
			return
		}

		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM user WHERE login = ?", req.Login)
		if err != nil {
			logger.Log("lvl", "error", "msg", "checking if user count", "err", err)
			api.WriteError(w, 500, err)
			return
		}

		if count == 0 {
			_, err := db.Exec("INSERT INTO user (login, password) VALUES (?, ?)", req.Login, req.Password)
			if err != nil {
				logger.Log("lvl", "error", "msg", "registering new user", "err", err)
				api.WriteError(w, 500, err)
				return
			}

			logger.Log("lvl", "info", "msg", "user registered", "login", req.Login)
		}

		err = db.Get(&count, "SELECT COUNT(*) FROM user WHERE login = ? AND password = ?", req.Login, req.Password)
		if err != nil {
			logger.Log("lvl", "error", "msg", "checking user password", "err", err)
			api.WriteError(w, 500, err)
			return
		}

		if count == 0 {
			logger.Log("lvl", "error", "msg", "user login failed", "login", req.Login)
			api.WriteError(w, 403, errors.New("invalid password"))
			return
		}

		t, err := token.SignHS256(token.Claims{
			"user": req.Login,
			"exp":  time.Now().Add(tokenTTL).Unix(),
		}, secret)
		if err != nil {
			logger.Log("lvl", "error", "msg", "unable to generate token", "err", err)
			api.WriteError(w, 500, err)
		}

		api.Write(w, map[string]interface{}{
			"token": t,
		})
	}
}
