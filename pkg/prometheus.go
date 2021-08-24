package pkg

import (
	"context"
	"fmt"
	"github.com/CloudForecast/barometer/pkg/barometerApi"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"time"
)

var prometheusURL string

func NewPrometheusAPIClient() v1.API {
	prometheusURL = viper.GetString("prometheusUrl")
	client, err := api.NewClient(api.Config{
		Address: prometheusURL,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	return v1.NewAPI(client)
}

func ExecutePromQLQuery(v1api v1.API, query string, start time.Time, end time.Time, duration time.Duration) (model.Value, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	result, warnings, err := v1api.QueryRange(ctx, query, v1.Range{
		Start: start,
		End:   end,
		Step:  duration,
	})
	if err != nil {
		log.Error().Msgf("unexpected error: %v\n", err)
		return nil, err
	}
	if len(warnings) > 0 {
		log.Warn().Msgf("warnings: %v\n", warnings)
	}
	log.Debug().Msgf("ran query '%s'", query)
	log.Trace().Msgf("ran query '%s', query result: %v", query, result)

	return result, nil
}

func followPromQlInstructions(instruction *barometerApi.PromQlQueryInstruction, resultsChan chan<- []barometerApi.PromQLResult, errorsChan chan<- error) {
	promClient := NewPrometheusAPIClient()

	// If needed, this could be parallelized later. For now, a simple loop should suffice.
	for _, config := range instruction.Configurations {
		var configurationResults []barometerApi.PromQLResult

		for queryName, query := range instruction.Queries {
			results, err := ExecutePromQLQuery(promClient, string(query), time.Unix(int64(config.StartTs), 0), time.Unix(int64(config.EndTs), 0), time.Duration(config.StepSec)*time.Second)
			if err != nil {
				log.Error().Err(err).Msg("")
				errorsChan <- err
				continue
			}
			var convertedResults []interface{}
			mapstructure.Decode(results, &convertedResults)
			if results != nil {
				log.Trace().Msgf("returned result: %v", results)
				var result = barometerApi.PromQLResult{
					QueryId:             string(queryName),
					Query:               string(query),
					PromQlConfiguration: config,
					Result:              convertedResults,
				}
				configurationResults = append(configurationResults, result)
			}
		}
		resultsChan <- configurationResults
	}
	close(resultsChan)
	close(errorsChan)
}

func FetchAndSubmitPrometheusData(b barometerApi.ApiClient) error {
	instructions, err := b.GetPromQlInstructions()
	if err != nil {
		return errors.Wrap(err, "issue getting promql instructions")
	}

	resultsChan := make(chan []barometerApi.PromQLResult)
	errorsChan := make(chan error)
	var errorList []error
	go followPromQlInstructions(instructions, resultsChan, errorsChan)

	for resultsChan != nil || errorsChan != nil {
		select {
		case results, ok := <-resultsChan:
			if !ok {
				resultsChan = nil
			} else {
				err = b.SendPromQlResultsEvent(*instructions, results)
				if err != nil {
					log.Error().Err(err).Msg("Error sending promql results event")
					errorList = append(errorList, err)
				}
			}

		case err, ok := <-errorsChan:
			if !ok {
				errorsChan = nil
			} else {
				errorList = append(errorList, err)
			}
		}
	}

	for _, err = range errorList {
		log.Error().Err(err).Msg("error following PromQl instructions")
		go b.SendExceptionEvent(err)
	}

	return nil
}
