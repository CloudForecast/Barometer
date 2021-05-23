package barometerApi

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

func NewExceptionEvent(err error) BarometerEvent {
	var outputData map[string]interface{}
	eventData := BarometerExceptionEventData{
		ErrorBool: true,
		Message:   err.Error(),
	}
	mapstructure.Decode(eventData, &outputData)

	return BarometerEvent{
		EventType: Exception,
		Event:     outputData,
	}
}

// SendExceptionEvent takes an error rather than event data, packages it up, and
// sends it the Barometer for logging.
func (b BarometerApi) SendExceptionEvent(inputError error) error {
	log.Debug().Msgf("sending exception event: %v", inputError)
	event := NewExceptionEvent(inputError)
	statusCode, err := b.makePostRequest(event)
	if err != nil {
		return err
	}
	if statusCode != 200 {
		return fmt.Errorf("received unexpected status code %d from sending exception", statusCode)
	}
	return nil
}
