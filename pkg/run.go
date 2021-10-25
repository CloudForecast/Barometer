package pkg

import (
	"github.com/CloudForecast/barometer/pkg/barometerApi"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"math/rand"
	"strconv"
	"time"
)

func RunAll(client barometerApi.ApiClient) (func(), error) {
	// Retrieve cronSchedule and generate new one if needed
	s := gocron.NewScheduler(time.UTC)
	cronSchedule := viper.GetString("schedule")
	if cronSchedule == "auto-generated" {
		// Generate a new cronSchedule to better spread the requests on the BM API
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		minutes := 15 + r1.Intn(30)
		cronSchedule = strconv.Itoa(minutes) + " * * * *"
	}

	appVersion := viper.GetString("appVersion")
	healthCheckInfo := map[string]interface{}{"appVersion": appVersion, "cronSchedule": cronSchedule}

	// Set the HealthCheck Cron Job
	stopHealthChecks, err := BeginHealthChecks(client, healthCheckInfo)
	if err != nil {
		return nil, err
	}

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
