package server

import (
	"github.com/gilcrest/httplog"
	"github.com/gilcrest/servertoken"
	"github.com/gilcrest/srvr/datastore"
	"github.com/justinas/alice"
)

// routes registers handlers to the router
func (s *Server) routes() error {

	// Get logging Database to pass into httplog
	// Only need this if you plan to use the PostgreSQL
	// logging style of httplog
	logdb, err := s.DS.DB(datastore.LogDB)
	if err != nil {
		return err
	}

	// Get App Database for token authentication
	appdb, err := s.DS.DB(datastore.AppDB)
	if err != nil {
		return err
	}

	// httplog.NewOpts gets a new httplog.Opts struct
	// (with all flags set to false)
	opts := new(httplog.Opts)

	// For the examples below, I chose to turn on db logging only
	// Log the request headers only (body has password on this api!)
	// Log both the response headers and body
	opts.Log2DB.Enable = true
	opts.Log2DB.Request.Header = true
	opts.Log2DB.Response.Header = true
	opts.Log2DB.Response.Body = true

	// match only POST requests on /api/v1/client
	// with Content-Type header = application/json
	s.Router.Handle("/v1/client",
		alice.New(httplog.LogHandler(s.Logger, logdb, opts), s.handleRespHeader, servertoken.Handler(s.Logger, appdb)).
			ThenFunc(s.handleClient())).
		Methods("POST").
		Headers("Content-Type", "application/json")

	return nil
}
