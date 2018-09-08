package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"

	"github.com/jasonrichardsmith/sentry/config"
	_ "github.com/jasonrichardsmith/sentry/healthz"
	_ "github.com/jasonrichardsmith/sentry/limits"
	_ "github.com/jasonrichardsmith/sentry/nslabels"
	_ "github.com/jasonrichardsmith/sentry/source"
	_ "github.com/jasonrichardsmith/sentry/tags"

	log "github.com/Sirupsen/logrus"
	"github.com/jasonrichardsmith/sentry/mux"
	"github.com/jasonrichardsmith/sentry/sentry"
)

var (
	dev bool
)

func init() {
	flag.BoolVar(&dev, "dev", false, "Run in dev mode no tls.")
}

func main() {
	log.Info("Loading Config")
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Config loaded")
	log.Info(config.DefaultConfig)
	s := mux.New(config.DefaultConfig)
	var ss *http.Server

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := ss.Shutdown(context.Background()); err != nil {
			log.Printf("Sentry server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if dev {
		log.Info("Serving Sentry without TLS")
		ss = sentry.NewSentryServerNoSSL(s)
		log.Fatal(ss.ListenAndServe())
	} else {
		log.Info("Starting new Sentry Server")
		ss, err = sentry.NewSentryServer(s)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(ss.ListenAndServeTLS("", ""))
	}
	<-idleConnsClosed
}
