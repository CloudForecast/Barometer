package pkg

import (
	"github.com/CloudForecast/Barometer/pkg/barometerApi"
	"time"
	"github.com/go-co-op/gocron"
)

func sendHealthCheck(b barometerApi.ApiClient) error {
	return b.SendHealthCheckEvent()
}

func BeginHealthChecks(b barometerApi.ApiClient) (func(), error) {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(5).Minutes().SingletonMode().Do(func() {
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