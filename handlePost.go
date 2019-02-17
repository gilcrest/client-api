package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gilcrest/apiclient"
	"github.com/gilcrest/errors"
	"github.com/gilcrest/httplog"
	"github.com/gilcrest/srvr/datastore"
)

// CreateClientHandler is used to create a new client (aka app)
// and generate clientID, clientSecret, etc.
func (s *Server) handleClient() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		// request is used for the POST /client API request
		type request struct {
			ClientName        string `json:"client_name"`
			ClientHomeURL     string `json:"homepage_url"`
			ClientDescription string `json:"client_description"`
			RedirectURI       string `json:"redirect_uri"`
			Username          string `json:"username"`
		}

		// response is used for the /client API responses
		type response struct {
			ClientID          string         `json:"client_id"`
			ClientName        string         `json:"client_name"`
			ClientHomeURL     string         `json:"homepage_url"`
			ClientDescription string         `json:"client_description"`
			RedirectURI       string         `json:"redirect_uri"`
			PrimaryUserID     string         `json:"username"`
			ClientSecret      string         `json:"client_secret"`
			ServerToken       string         `json:"server_token"`
			CreateTimestamp   string         `json:"create_timestamp"`
			Audit             *httplog.Audit `json:"audit"`
		}

		// retrieve the context from the http.Request
		ctx := req.Context()

		rqst := new(request)
		err := json.NewDecoder(req.Body).Decode(&rqst)
		defer req.Body.Close()
		if err != nil {
			err = errors.RE(http.StatusBadRequest, errors.InvalidRequest, err)
			errors.HTTPError(w, err)
			return
		}

		client := new(apiclient.Client)

		client.Name = rqst.ClientName
		client.HomeURL = rqst.ClientHomeURL
		client.Description = rqst.ClientDescription
		client.RedirectURI = rqst.RedirectURI
		client.PrimaryUserID = rqst.Username

		err = client.Finalize()
		if err != nil {
			err = errors.RE(http.StatusBadRequest, errors.Validation, err)
			errors.HTTPError(w, err)
			return
		}

		// get a new DB Tx
		tx, err := s.DS.BeginTx(ctx, nil, datastore.AppDB)
		if err != nil {
			err = errors.RE(http.StatusInternalServerError, errors.Database)
			errors.HTTPError(w, err)
			return
		}

		// Call the CreateClientDB method of the Client object
		// to write to the db
		err = client.CreateClientDB(ctx, tx)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.RE(http.StatusInternalServerError, errors.Database)
				errors.HTTPError(w, err)
				return
			}
		}

		if err := tx.Commit(); err != nil {
			err = errors.RE(http.StatusInternalServerError, errors.Database)
			errors.HTTPError(w, err)
			return
		}

		// If we successfully committed the db transaction, we can consider this
		// transaction successful and return a response with the response body

		aud, err := httplog.NewAudit(ctx)
		if err != nil {
			err = errors.RE(http.StatusInternalServerError, errors.Unanticipated)
			errors.HTTPError(w, err)
			return
		}

		// create a new response struct and set Audit and other
		// relevant elements
		resp := new(response)
		resp.Audit = aud
		resp.ClientID = client.ID
		resp.ClientName = client.Name
		resp.ClientHomeURL = client.HomeURL
		resp.ClientDescription = client.Description
		resp.RedirectURI = client.RedirectURI
		resp.PrimaryUserID = client.PrimaryUserID
		resp.ClientSecret = client.Secret
		resp.ServerToken = client.ServerToken
		resp.CreateTimestamp = client.CreateTimestamp.Format(time.RFC3339)

		// Encode response struct to JSON for the response body
		json.NewEncoder(w).Encode(*resp)
		if err != nil {
			err = errors.RE(http.StatusInternalServerError, errors.Internal)
			errors.HTTPError(w, err)
			return
		}

	}
}
