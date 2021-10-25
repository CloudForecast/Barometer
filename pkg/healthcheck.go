package pkg

import (
	"github.com/CloudForecast/barometer/pkg/barometerApi"
	"github.com/go-co-op/gocron"
	"time"
)

// sendHealthCheck sends a healthcheck event to the Barometer API server
func sendHealthCheck(b barometerApi.ApiClient, info map[string]interface{}) error {
	return b.SendHealthCheckEvent(info)
}

func BeginHealthChecks(b barometerApi.ApiClient, info map[string]interface{}) (func(), error) {
	// Run first Health check
	_ = sendHealthCheck(b, info)

	// Setup cron job every 15 minutes
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(15).Minutes().SingletonMode().Do(func() {
		_ = sendHealthCheck(b, info)
	})
	if err != nil {
		return nil, err
	}
	s.StartAsync()
	return func() {
		s.Stop()
	}, nil
}
