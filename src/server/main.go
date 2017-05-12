package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/negroni"
)

type Configuration struct {
	Bind string
}

func main() {
	log := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))

	var c Configuration
	err := envconfig.Process("resistance", &c)
	if err != nil {
		log.Log("lvl", "error", "msg", "parsing configuration", "err", err)
		os.Exit(1)
	}
	log.Log("lvl", "info", "msg", "server starting")

	r := httprouter.New()
	r.GET("/", httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		write(w, map[string]interface{}{
			"name": "resistance",
		})
	}))

	n := negroni.New()
	n.UseHandler(r)

	s := &http.Server{
		Addr:    c.Bind,
		Handler: n,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Log("lvl", "error", "msg", "listen error", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.Shutdown(ctx)

	log.Log("lvl", "info", "msg", "server stopping")
}
