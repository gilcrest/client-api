package server

import (
	"net/http"

	"github.com/gilcrest/env"
	"github.com/gilcrest/env/datastore"
	"github.com/gilcrest/errors"
	"github.com/rs/zerolog"
)

// Server struct contains the environment (env.Env) and additional methods
// for running our HTTP server
type Server struct {
	*env.Env
}

// NewServer is a constructor for the Server struct
// Sets up the struct and registers routes
func NewServer(name env.Name, lvl zerolog.Level) (*Server, error) {
	const op errors.Op = "client-api/server/NewServer"

	env, err := env.NewEnv(name, lvl)
	if err != nil {
		return nil, errors.E(op, err)
	}

	server := &Server{env}

	err = server.DS.Option(datastore.InitLogDB())
	if err != nil {
		return nil, errors.E(op, err)
	}

	err = server.routes()
	if err != nil {
		return nil, errors.E(op, err)
	}

	return server, nil
}

// handleRespHeader middleware is used to add standard HTTP response headers
func (s *Server) handleRespHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			h.ServeHTTP(w, r) // call original
		})
}
