package barometerApi

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"time"
)

// NewPromQlResultsEvent creates a new event for a single PromQLResult to send to the Barometer API.
// As we may receive instructions to query the same metric over varying time periods in the same
// instruction set, and each would share a query ID like `cpu_requested`, but each API call to Barometer
// can only have the query_id once, we are only able to send a single PromQlResult at a time.
func NewPromQlResultsEvent(instructions PromQlQueryInstruction, results []PromQLResult) *BarometerEvent {
	// Sometimes we end up with an empty result at the end, so we need to filter it out
	// before sending it to the API.
	var filteredResults []PromQLResult
	for _, result := range results {
		if result.Query != "" {
			filteredResults = append(filteredResults, result)
		}
	}

	var outputMap map[string]interface{}
	err := mapstructure.Decode(BarometerPromQlResultsEventData{
		PromQLInstructions: instructions,
		PromQlResults: filteredResults,
	}, &outputMap)
	if err != nil {
		panic(err)
	}

	return &BarometerEvent{
		EventType: PromQlResults,
		EventTs:   time.Now().Unix(),
		Event: outputMap,
	}
}

func (b BarometerApi) SendPromQlResults(instructions PromQlQueryInstruction, promQLResultsWrapper PromQLResultsWrapper) error {
	log.Info().Msg("Sending PromQlResult data...")
	event := NewPromQlResultsEvent(instructions, promQLResultsWrapper.Results)

	// Upload data to S3
	statusCode, err := b.UploadDataToS3(promQLResultsWrapper.PromQlConfiguration.UploadConfiguration, event)
	if err != nil {
		return err
	}
	if statusCode != 200 {
		return fmt.Errorf("received unexpected status code %d from SendPromQlResults when sending the data to S3", statusCode)
	}

	err = b.SendUploadedDataEvent(promQLResultsWrapper.PromQlConfiguration.UploadConfiguration.EventUUID, PromQlResults)
	if err != nil {
		return err
	}
	return nil
}
