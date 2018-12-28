package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gilcrest/client-api/server"
	"github.com/gilcrest/errors"
	"github.com/rs/zerolog"
)

func main() {
	const op errors.Op = "srvr.main"

	// flag allows for setting logging level, e.g. to run the server
	// with level set to debug, it'd be: ./server loglvl=debug
	loglvlFlag := flag.String("loglvl", "error", "sets log level")

	flag.Parse()

	loglevel := logLevel(loglvlFlag)

	server, err := server.NewServer(loglevel)
	if err != nil {
		log.Fatal(err)
	}

	// handle all requests with the Gorilla router
	http.Handle("/", server.Router)

	// ListenAndServe on port 8008, not specifying a particular IP address
	// for this particular implementation
	if err := http.ListenAndServe(":8008", nil); err != nil {
		log.Fatal(err)
	}
}

func logLevel(s *string) zerolog.Level {

	var lvl zerolog.Level

	// dereference the string pointer to get flag value
	ds := *s

	switch ds {
	case "debug":
		lvl = zerolog.DebugLevel
	case "info":
		lvl = zerolog.InfoLevel
	case "warn":
		lvl = zerolog.WarnLevel
	case "fatal":
		lvl = zerolog.FatalLevel
	case "panic":
		lvl = zerolog.PanicLevel
	case "disabled":
		lvl = zerolog.Disabled
	default:
		lvl = zerolog.ErrorLevel
	}
	return lvl
}
