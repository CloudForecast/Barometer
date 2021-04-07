package barometerApi

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

func NewHealthCheckEvent() BarometerEvent {
	return BarometerEvent{
		EventKey: HealthCheck,
		Event: make(map[string]interface{}),
	}
}

func (b BarometerApi) SendHealthCheckEvent() error {
	log.Debug().Msg("Sending health check...")
	event := NewHealthCheckEvent()
	statusCode, err := b.makePostRequest(event)
	if err != nil {
		return err
	}
	if statusCode != 200 {
		return fmt.Errorf("received unexpected status code %d from healthhcheck", statusCode)
	}
	return nil
}