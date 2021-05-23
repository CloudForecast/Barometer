package barometerApi

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

// NewPromQlResultsEvent creates a new event for a single PromQLResult to send to the Barometer API.
// As we may receive instructions to query the same metric over varying time periods in the same
// instruction set, and each would share a query ID like `cpu_requested`, but each API call to Barometer
// can only have the query_id once, we are only able to send a single PromQlResult at a time.
func NewPromQlResultsEvent(result PromQLResult) *BarometerPromQlResultsEventData {
	return &BarometerPromQlResultsEventData{
		PromQlResults: map[string]PromQLResult{
			result.QueryId: result,
		},
	}
}

func (b BarometerApi) SendPromQlResultsEvent(eventData PromQLResult) error {
	log.Debug().Msg("Sending PromQlResult data...")
	event := NewPromQlResultsEvent(eventData)
	statusCode, err := b.makePostRequest(event)
	if err != nil {
		return err
	}
	if statusCode != 200 {
		return fmt.Errorf("received unexpected status code %d from prometheus API data submission", statusCode)
	}
	return nil
}
