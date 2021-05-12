package pkg

import (
	"github.com/CloudForecast/Barometer/pkg/barometerApi"
	"github.com/go-co-op/gocron"
	"time"
)

func sendHealthCheck(b barometerApi.ApiClient) error {
	return b.SendHealthCheckEvent()
}

func BeginHealthChecks(b barometerApi.ApiClient) (func(), error) {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Minutes().SingletonMode().Do(func() {
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