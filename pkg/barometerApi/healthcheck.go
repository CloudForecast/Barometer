package barometerApi

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

func NewHealthCheckEvent() BarometerEvent {

	appVersion := viper.GetString("appVersion")
	return BarometerEvent{
		EventType: HealthCheck,
		EventTs:   time.Now().Unix(),
		Event:     map[string]interface{}{"appVersion": appVersion},
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
