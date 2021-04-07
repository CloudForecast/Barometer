package pkg

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"time"
)

var prometheusURL string

func init() {
	prometheusURL = viper.GetString("prometheusUrl")
}

func NewPrometheusAPIClient() v1.API {
	client, err := api.NewClient(api.Config{
		Address: prometheusURL,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	return v1.NewAPI(client)
}

func ExecutePromQLQuery(v1api v1.API, query string, start time.Time, end time.Time, duration time.Duration) model.Value {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.QueryRange(ctx, query, v1.Range{
		Start: start,
		End: end,
		Step: duration,
	})
	if err != nil {
		log.Error().Msgf("unexpected error: %v\n", err)
	}
	if len(warnings) > 0 {
		log.Warn().Msgf("warnings: %v\n", warnings)
	}
	log.Debug().Msgf("ran query '%s', query result: %v", query, result)

	return result
}
