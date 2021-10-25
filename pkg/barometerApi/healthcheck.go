package barometerApi

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

func NewHealthCheckEvent(info map[string]interface{}) BarometerEvent {
	return BarometerEvent{
		EventType: HealthCheck,
		EventTs:   time.Now().Unix(),
		Event:     info,
	}
}

func (b BarometerApi) SendHealthCheckEvent(info map[string]interface{}) error {
	log.Debug().Msg("Sending health check...")
	event := NewHealthCheckEvent(info)
	statusCode, err := b.makePostRequest(event)
	if err != nil {
		return err
	}
	if statusCode != 200 {
		return fmt.Errorf("received unexpected status code %d from healthhcheck", statusCode)
	}
	return nil
}


