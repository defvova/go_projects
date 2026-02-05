package writejson

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func NewJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if v == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Err(err).Msg("writeJSON encode error")
	}
}
