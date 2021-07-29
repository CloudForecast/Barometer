package pkg

import (
	"github.com/CloudForecast/Barometer/pkg/barometerApi"
	"github.com/go-co-op/gocron"
	"time"
)

// sendHealthCheck sends a healthcheck event to the Barometer API server
func sendHealthCheck(b barometerApi.ApiClient) error {
	return b.SendHealthCheckEvent()
}

func BeginHealthChecks(b barometerApi.ApiClient) (func(), error) {
	// Run first Health check
	_ = sendHealthCheck(b)

	// Setup cron job every 15 minutes
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(15).Minutes().SingletonMode().Do(func() {
		_ = sendHealthCheck(b)
	})
	if err != nil {
		return nil, err
	}
	s.StartAsync()
	return func() {
		s.Stop()
	}, nil
}
