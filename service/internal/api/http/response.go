package http

import (
	"encoding/json"
	"net/http"
)

func respond(w http.ResponseWriter, data []byte, statusCode int) error {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func (s server) respondError(w http.ResponseWriter, err error, code int) {
	s.logger.Error().Err(err).Msg("error responding to request")

	bytes, _ := json.Marshal(map[string]string{"error": err.Error()})

	if err := respond(w, bytes, code); err != nil {
		s.logger.Error().Err(err).Msg("error responding to request")

	}
}
