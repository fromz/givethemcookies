package main

import (
	"github.com/fromz/givethemcookies"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var log *logrus.Logger

func main() {
	log = logrus.New()
	givethemcookies.SetLogger(log)
	JWTVerifier, err := givethemcookies.GetJWTVerifierFromENVVars()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// set up http server with jwt verifier as mddleware
	http.HandleFunc("/", givethemcookies.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {

	}, JWTVerifier, claimFetcher))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func claimFetcher() interface{} {
	return givethemcookies.Claims{}
}
