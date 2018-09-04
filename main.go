package main

import (
	"flag"
	"net/http"

	"github.com/jasonrichardsmith/sentry/config"
	_ "github.com/jasonrichardsmith/sentry/limits"

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
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	s := mux.New(config.DefaultConfig)
	var ss *http.Server
	if dev {
		ss = sentry.NewSentryServerNoSSL(s)
		log.Fatal(ss.ListenAndServe())
	} else {
		ss, err = sentry.NewSentryServer(s)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(ss.ListenAndServeTLS("", ""))
	}
}
