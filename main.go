package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/jasonrichardsmith/sentry/mux"
	"github.com/jasonrichardsmith/sentry/sentry"
	"net/http"
)

var (
	dev bool
)

func init() {
	flag.BoolVar(&dev, "dev", false, "Run in dev mode no tls.")
}

func main() {
	config := mux.New()
	s, err := config.LoadSentry()
	if err != nil {
		log.Fatal(err)
	}
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
