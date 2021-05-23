package pkg

import (
	"context"
	"fmt"
	"github.com/CloudForecast/Barometer/pkg/barometerApi"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
	log.Debug().Msgf("ran query '%s', query result: %v", query, result)

	return result, nil
}

func followPromQlInstructions(instruction *barometerApi.PromQlQueryInstruction, resultsChan chan<- barometerApi.PromQLResult, errorList *[]error) {
	promClient := NewPrometheusAPIClient()
	var localErrorList []error

	// If needed, this could be parallelized later. For now, a simple loop should suffice.
	for _, config := range instruction.Configurations {
		for queryName, query := range instruction.Queries {
			results, err := ExecutePromQLQuery(promClient, string(query), time.Unix(int64(config.StartTs), 0), time.Unix(int64(config.EndTs), 0), time.Duration(config.StepSec)*time.Second)
			if err != nil {
				log.Error().Err(err).Msg("")
				localErrorList = append(localErrorList, err)
				continue
			}
			var convertedResults []interface{}
			mapstructure.Decode(results, &convertedResults)
			if results != nil {
				log.Debug().Msgf("returned result: %v", results)
				resultsChan <- barometerApi.PromQLResult{
					QueryId:             string(queryName),
					Query:               string(query),
					PromQlConfiguration: config,
					Result:              convertedResults,
				}
			}
		}
	}

	close(resultsChan)
	errorList = &localErrorList
}

func FetchAndSubmitPrometheusData(b barometerApi.ApiClient) error {
	instructions, err := b.GetPromQlInstructions()
	if err != nil {
		return errors.Wrap(err, "issue getting promql instructions")
	}

	resultsChan := make(chan barometerApi.PromQLResult)
	var errorList []error
	go followPromQlInstructions(instructions, resultsChan, &errorList)

	for {
		result, more := <-resultsChan
		if !more {
			break
		}
		go func() {
			err := b.SendPromQlResultsEvent(result)
			if err != nil {
				log.Error().Err(err).Msg("Error sending promql results event")
				_ = b.SendExceptionEvent(err)
			}
		}()
	}

	for _, err = range errorList {
		log.Error().Err(err).Msg("error following PromQl instructions")
		go b.SendExceptionEvent(err)
	}

	return nil
}
