// Package http provides the http server and handlers.
package http

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	v1 "github.com/mnabbasabadi/relex_convertor/api/rest/v1"
	"github.com/mnabbasabadi/relex_convertor/service/internal/logic"
	"github.com/rs/zerolog"
)

const (
	contentTypeCSV = "text/csv"
)

var _ v1.ServerInterface = new(server)

type server struct {
	logger zerolog.Logger
	logic  logic.Logic
}

// Convert ...
func (s server) Convert(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != contentTypeCSV {
		s.respondError(w, fmt.Errorf("expected text/csv, got %s", contentType), http.StatusBadRequest)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			s.logger.Error().Err(err).Msg("Failed to close body")
		}
	}()

	reader := csv.NewReader(r.Body)
	// read all the records
	records, err := reader.ReadAll()
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to read all records")
		s.respondError(w, err, http.StatusBadRequest)
		return
	}

	node, err := s.logic.ConvertToTreeNode(r.Context(), records)
	if err != nil {
		if errors.Is(err, logic.ErrorInvalidPayload) {
			s.logger.Error().Err(err).Msg("Failed to convert to JSON")
			s.respondError(w, err, http.StatusBadRequest)
			return
		}
		s.logger.Error().Err(err).Msg("Failed to convert to JSON")
		s.respondError(w, err, http.StatusInternalServerError)

	}
	responseBody, err := json.Marshal(node)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to generate JSON")
		s.respondError(w, err, http.StatusInternalServerError)
	}

	if err := respond(w, responseBody, http.StatusOK); err != nil {
		s.logger.Error().Err(err).Msg("Failed to respond")
		s.respondError(w, err, http.StatusInternalServerError)
	}

}

// NewServer returns a new http.Handler that implements the api.ServerInterface
func NewServer(logic logic.Logic, logger zerolog.Logger) http.Handler {
	s := server{
		logic:  logic,
		logger: logger,
	}

	options := v1.ChiServerOptions{
		BaseRouter: chi.NewRouter(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Error().Err(err).Msg("error")
			s.respondError(w, err, http.StatusBadRequest)
		},
	}

	return v1.HandlerWithOptions(s, options)

}
