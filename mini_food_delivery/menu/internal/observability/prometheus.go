package observability

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func StartPrometheusServer(addr string) {
	log.Info().Msgf("prometheus metrics on %s", addr)
	if err := http.ListenAndServe(addr, promhttp.Handler()); err != nil {
		log.Fatal().Err(err).Msg("prometheus server failed")
	}
}
