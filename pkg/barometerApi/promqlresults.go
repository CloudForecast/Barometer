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

func (b BarometerApi) SendPromQlResultsEvent(instructions PromQlQueryInstruction, eventData []PromQLResult) error {
	log.Info().Msg("Sending PromQlResult data...")
	event := NewPromQlResultsEvent(instructions, eventData)
	statusCode, err := b.makePostRequest(event)
	if err != nil {
		return err
	}
	if statusCode != 200 {
		return fmt.Errorf("received unexpected status code %d from prometheus API data submission", statusCode)
	}
	return nil
}
