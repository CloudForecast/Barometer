package pkg

import (
	"github.com/CloudForecast/Barometer/pkg/barometerApi"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"time"
)

func RunAll(client barometerApi.ApiClient) (func(), error) {
	// Health checks
	stopHealthChecks, err := BeginHealthChecks(client)
	if err != nil {
		return nil, err
	}

	// Everything else
	s := gocron.NewScheduler(time.UTC)
	_, err = s.Every(5).Minutes().SingletonMode().Do(func() {
		go func() {
			err := FetchAndSubmitKubernetesObjects(client)
			if err != nil {
				log.Error().Err(err).Msg("")
			}
		}()

		go func() {
			err = FetchAndSubmitPrometheusData(client)
			if err != nil {
				log.Error().Err(err).Msg("")
			}
		}()
	})
	s.StartAsync()

	if err != nil {
		return nil, err
	}

	return func() {
		s.Stop()
		stopHealthChecks()
	}, nil
}
