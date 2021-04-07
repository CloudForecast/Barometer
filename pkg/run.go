package pkg

import "github.com/CloudForecast/Barometer/pkg/barometerApi"

func RunAll(client barometerApi.ApiClient) (func(), error) {
	// Health checks
	stopHealthChecks, err := BeginHealthChecks(client)
	if err != nil {
		return nil, err
	}

	return func() {
		stopHealthChecks()
	}, nil
}