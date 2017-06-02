package main

import (
	"api"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

type Configuration struct {
	Bind     string
	Database string
	Secret   string
	TokenTTL time.Duration
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))

	var c Configuration
	err := envconfig.Process("resistance", &c)
	if err != nil {
		logger.Log("lvl", "error", "msg", "parsing configuration", "err", err)
		os.Exit(1)
	}
	logger.Log("lvl", "info", "msg", "server starting")

	db, err := sqlx.Connect("sqlite3", c.Database)
	if err != nil {
		logger.Log("lvl", "error", "msg", "connecting to database", "err", err)
		os.Exit(1)
	}
	logger.Log("lvl", "info", "msg", "connected to database")

	r := httprouter.New()
	r.GET("/", httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		api.Write(w, map[string]interface{}{
			"name": "resistance",
		})
	}))
	r.POST("/login", LoginHandler(db, []byte(c.Secret), c.TokenTTL))
	r.POST("/authenticate", AuthenticateHandler(db, []byte(c.Secret)))
	r.GET("/game", ListGameHandler(db))
	r.POST("/game", CreateGameHandler(db))
	r.GET("/game/:id", ShowGameHandler(db))

	n := negroni.New()
	n.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
	n.Use(NewContextMiddleware(Context{
		Logger: logger,
	}))
	n.Use(AuthenticateMiddleware(logger, db, []byte(c.Secret), []string{
		"/login",
		"/authenticate",
	}))
	n.UseHandler(r)

	s := &http.Server{
		Addr:    c.Bind,
		Handler: n,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			logger.Log("lvl", "error", "msg", "listen error", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.Shutdown(ctx)

	logger.Log("lvl", "info", "msg", "server stopping")
}
