package main

import (
	"api"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/elwinar/token"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

const (
	ClaimUser = "user"
	ClaimExp  = "exp"
)

func LoginHandler(db *sqlx.DB, secret []byte, tokenTTL time.Duration) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ctx := Ctx(r)
		logger := log.With(ctx.Logger, "handler", "login")

		var req LoginRequest
		err := api.Read(r, &req)
		if err != nil {
			logger.Log("lvl", "error", "msg", "reading request", "err", err.Error())
			api.WriteError(w, http.StatusBadRequest, err)
			return
		}

		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM user WHERE login = ?", req.Login)
		if err != nil {
			logger.Log("lvl", "error", "msg", "checking if user count", "err", err.Error())
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if count == 0 {
			hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
			_, err := db.Exec("INSERT INTO user (login, password) VALUES (?, ?)", req.Login, string(hash))
			if err != nil {
				logger.Log("lvl", "error", "msg", "registering new user", "err", err.Error())
				api.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			logger.Log("lvl", "info", "msg", "user registered", "login", req.Login)
		}

		var user User
		err = db.Get(&user, "SELECT id, login, password FROM user WHERE login = ?", req.Login)
		if err != nil {
			logger.Log("lvl", "error", "msg", "retrieving user", "err", err.Error())
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
			err = errors.New("invalid password")
			logger.Log("lvl", "error", "msg", "login failed", "err", err.Error())
			api.WriteError(w, http.StatusForbidden, err)
			return
		}

		t, err := token.SignHS256(token.Claims{
			ClaimUser: user.ID,
			ClaimExp:  time.Now().Add(tokenTTL).Unix(),
		}, secret)
		if err != nil {
			logger.Log("lvl", "error", "msg", "unable to generate token", "err", err.Error())
			api.WriteError(w, http.StatusInternalServerError, err)
		}

		api.Write(w, map[string]interface{}{
			"token": t,
		})
	}
}
