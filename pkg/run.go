package pkg

import (
	"github.com/CloudForecast/barometer/pkg/barometerApi"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

func RunAll(client barometerApi.ApiClient) (func(), error) {
	// Set the HealthCheck Cron Job
	stopHealthChecks, err := BeginHealthChecks(client)
	if err != nil {
		return nil, err
	}

	// Everything else
	s := gocron.NewScheduler(time.UTC)
	cronSchedule := viper.GetString("schedule")

	// Trigger the first RunFetchAndSubmitKubernetesObjectsAndPrometheusData(client)
	RunFetchAndSubmitKubernetesObjectsAndPrometheusData(client)

	log.Error().Msgf("cronSchedule: %s", cronSchedule)
	_, err = s.Cron(cronSchedule).SingletonMode().Do(func() { RunFetchAndSubmitKubernetesObjectsAndPrometheusData(client) })
	s.StartAsync()

	if err != nil {
		return nil, err
	}

	return func() {
		s.Stop()
		stopHealthChecks()
	}, nil
}

func RunFetchAndSubmitKubernetesObjectsAndPrometheusData(client barometerApi.ApiClient) {
	log.Info().Msg("Triggering RunFetchAndSubmitKubernetesObjectsAndPrometheusData")
	go func() {
		err := FetchAndSubmitKubernetesObjects(client)
		if err != nil {
			log.Error().Err(err).Msg("")
		}
	}()

	go func() {
		err := FetchAndSubmitPrometheusData(client)
		if err != nil {
			log.Error().Err(err).Msg("")
		}
	}()
}
